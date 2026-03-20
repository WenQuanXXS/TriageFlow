package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/triageflow/backend/model"
	"github.com/triageflow/backend/service"
	"gorm.io/gorm"
)

type TaskHandler struct {
	DB            *gorm.DB
	TriageService service.Triager
	RuleEngine    *service.RuleEngine
}

type CreateTaskRequest struct {
	PatientName    string `json:"patient_name" binding:"required"`
	ChiefComplaint string `json:"chief_complaint" binding:"required"`
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
		Status:         "pending",
		Priority:       "normal",
		TriageStatus:   "pending",
	}

	// Run triage if service is available
	if h.TriageService != nil {
		triageResult, rawOutput, err := h.TriageService.PerformTriage(context.Background(), req.ChiefComplaint)
		if err != nil {
			task.TriageStatus = "failed"
			task.LLMRawOutput = err.Error()
		} else {
			symptomsJSON, _ := json.Marshal(triageResult.Symptoms)
			riskJSON, _ := json.Marshal(triageResult.RiskSignals)
			deptsJSON, _ := json.Marshal(triageResult.CandidateDepts)

			task.Symptoms = string(symptomsJSON)
			task.RiskSignals = string(riskJSON)
			task.CandidateDepts = string(deptsJSON)
			task.AISuggestedPri = triageResult.SuggestedPri
			task.LLMRawOutput = rawOutput
			task.TriageStatus = "completed"

			// Run rule engine
			if h.RuleEngine != nil {
				ruleResult := h.RuleEngine.Evaluate(req.ChiefComplaint, triageResult)
				task.RuleTriggered = ruleResult.RuleTriggered
				task.RuleReason = ruleResult.Reason
				task.FinalPriority = ruleResult.FinalPriority
				task.FinalDepartment = ruleResult.FinalDepartment
				task.Priority = ruleResult.FinalPriority
				task.Department = ruleResult.FinalDepartment
			}
		}
	}

	if err := h.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	var task model.Task
	if err := h.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	var tasks []model.Task
	query := h.DB

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if priority := c.Query("priority"); priority != "" {
		query = query.Where("final_priority = ? OR priority = ?", priority, priority)
	}

	if err := query.Order("created_at desc").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return empty array instead of null
	if tasks == nil {
		tasks = []model.Task{}
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "task already completed"})
		return
	}

	if err := h.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}
