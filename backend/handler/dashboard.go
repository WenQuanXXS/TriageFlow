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
	Total      int64           `json:"total"`
	ByStatus   []StatusCount   `json:"by_status"`
	ByPriority []PriorityCount `json:"by_priority"`
}

func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	var resp DashboardResponse

	h.DB.Table("tasks").Count(&resp.Total)
	h.DB.Table("tasks").Select("status, count(*) as count").Group("status").Scan(&resp.ByStatus)
	h.DB.Table("tasks").Select("priority, count(*) as count").Group("priority").Scan(&resp.ByPriority)

	c.JSON(http.StatusOK, resp)
}
