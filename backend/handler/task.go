package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/triageflow/backend/model"
	"gorm.io/gorm"
)

type TaskHandler struct {
	DB *gorm.DB
}

type CreateTaskRequest struct {
	PatientName    string `json:"patient_name" binding:"required"`
	ChiefComplaint string `json:"chief_complaint" binding:"required"`
	Priority       string `json:"priority"`
	Department     string `json:"department"`
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := model.Task{
		PatientName:    req.PatientName,
		ChiefComplaint: req.ChiefComplaint,
		Priority:       req.Priority,
		Department:     req.Department,
	}
	if task.Priority == "" {
		task.Priority = "normal"
	}
	if task.Status == "" {
		task.Status = "pending"
	}

	if err := h.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	var tasks []model.Task
	query := h.DB

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
	}

	if err := query.Order("created_at desc").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) ToggleStatus(c *gin.Context) {
	id := c.Param("id")
	var task model.Task
	if err := h.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	switch task.Status {
	case "pending":
		task.Status = "in_progress"
	case "in_progress":
		task.Status = "completed"
	case "completed":
		task.Status = "pending"
	}

	if err := h.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}
