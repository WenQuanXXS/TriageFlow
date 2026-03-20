package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	DB *gorm.DB
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type PriorityCount struct {
	Priority string `json:"priority"`
	Count    int64  `json:"count"`
}

type DashboardResponse struct {
	Total         int64           `json:"total"`
	ByStatus      []StatusCount   `json:"by_status"`
	ByPriority    []PriorityCount `json:"by_priority"`
	TriageCount   int64           `json:"triage_count"`
	RuleOverrides int64           `json:"rule_overrides"`
}

func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	var resp DashboardResponse

	if err := h.DB.Table("tasks").Count(&resp.Total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Table("tasks").Select("status, count(*) as count").Group("status").Scan(&resp.ByStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Table("tasks").Select("priority, count(*) as count").Group("priority").Scan(&resp.ByPriority).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Table("tasks").Where("triage_status = ?", "completed").Count(&resp.TriageCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.DB.Table("tasks").Where("rule_triggered != '' AND rule_triggered IS NOT NULL").Count(&resp.RuleOverrides).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
