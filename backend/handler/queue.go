package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/triageflow/backend/service"
	"gorm.io/gorm"
)

type QueueHandler struct {
	DB           *gorm.DB
	QueueService *service.QueueService
}

// ListQueue returns queue entries in backend-defined priority order.
// Supports filtering by department and queue_status.
func (h *QueueHandler) ListQueue(c *gin.Context) {
	dept := c.Query("department")
	status := c.Query("queue_status")

	entries, err := h.QueueService.ListByDepartment(h.DB, dept, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}

// GetPosition returns the patient's queue position, ahead count,
// plus the department's now-serving and waiting lists.
func (h *QueueHandler) GetPosition(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	pos, err := h.QueueService.GetPosition(h.DB, uint(taskID))
	if err != nil {
		if errors.Is(err, service.ErrQueueEntryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "queue entry not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, pos)
}

// CallPatient transitions a queue entry from "waiting" to "called"
// and updates the associated task status to "in_progress".
func (h *QueueHandler) CallPatient(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	entry, err := h.QueueService.CallPatient(h.DB, uint(taskID))
	if err != nil {
		if errors.Is(err, service.ErrQueueEntryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "queue entry not found"})
		} else if errors.Is(err, service.ErrInvalidTransition) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "patient is not in waiting status"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, entry)
}

// CompletePatient transitions a queue entry from "called" to "completed"
// and updates the associated task status to "completed".
func (h *QueueHandler) CompletePatient(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	entry, err := h.QueueService.CompletePatient(h.DB, uint(taskID))
	if err != nil {
		if errors.Is(err, service.ErrQueueEntryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "queue entry not found"})
		} else if errors.Is(err, service.ErrInvalidTransition) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "patient is not in called status"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, entry)
}
