package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/triageflow/backend/model"
	"github.com/triageflow/backend/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	host := getEnvOrDefault("DB_HOST", "127.0.0.1")
	port := getEnvOrDefault("DB_PORT", "3306")
	user := getEnvOrDefault("DB_USER", "root")
	pass := getEnvOrDefault("DB_PASS", "1234")
	name := getEnvOrDefault("DB_NAME", "triageflow_test")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("skipping: cannot connect to test database: %v", err)
	}

	db.Migrator().DropTable(&model.Task{})
	db.AutoMigrate(&model.Task{})
	db.Migrator().DropTable(&model.QueueEntry{})
	db.AutoMigrate(&model.QueueEntry{})

	return db
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	taskHandler := &TaskHandler{
		DB:            db,
		TriageService: service.NewMockTriageService(),
		RuleEngine:    service.NewRuleEngine(),
		QueueService:  service.NewQueueService(),
	}
	dashHandler := &DashboardHandler{DB: db}

	api := r.Group("/api")
	{
		api.POST("/tasks", taskHandler.CreateTask)
		api.GET("/tasks", taskHandler.ListTasks)
		api.GET("/tasks/:id", taskHandler.GetTask)
		api.PATCH("/tasks/:id/status", taskHandler.ToggleStatus)
		api.GET("/dashboard", dashHandler.GetDashboard)
	}

	return r
}

func TestCreateTask_Normal(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	body := `{"patient_name":"张三","chief_complaint":"头痛三天，伴有轻微恶心"}`
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var task model.Task
	json.Unmarshal(w.Body.Bytes(), &task)

	if task.PatientName != "张三" {
		t.Errorf("expected patient_name '张三', got '%s'", task.PatientName)
	}
	if task.TriageStatus != "completed" {
		t.Errorf("expected triage_status 'completed', got '%s'", task.TriageStatus)
	}
	if task.RuleTriggered != "" {
		t.Errorf("expected no rule triggered, got '%s'", task.RuleTriggered)
	}
	if task.FinalPriority == "" {
		t.Error("expected final_priority to be set")
	}
}

func TestCreateTask_HighRisk(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	body := `{"patient_name":"李四","chief_complaint":"突发胸痛，持续半小时，出汗"}`
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var task model.Task
	json.Unmarshal(w.Body.Bytes(), &task)

	if task.RuleTriggered != "chest_pain" {
		t.Errorf("expected rule 'chest_pain', got '%s'", task.RuleTriggered)
	}
	if task.FinalPriority != "urgent" {
		t.Errorf("expected final_priority 'urgent', got '%s'", task.FinalPriority)
	}
	if task.FinalDepartment != "Emergency" {
		t.Errorf("expected final_department 'Emergency', got '%s'", task.FinalDepartment)
	}
}

func TestGetTask(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	task := model.Task{PatientName: "Test", ChiefComplaint: "test", Status: "pending", Priority: "normal"}
	db.Create(&task)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/%d", task.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var got model.Task
	json.Unmarshal(w.Body.Bytes(), &got)
	if got.ID != task.ID {
		t.Errorf("expected id %d, got %d", task.ID, got.ID)
	}
}

func TestListTasks(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	db.Create(&model.Task{PatientName: "A", ChiefComplaint: "cough", Status: "pending", Priority: "normal"})
	db.Create(&model.Task{PatientName: "B", ChiefComplaint: "fever", Status: "in_progress", Priority: "urgent"})
	db.Create(&model.Task{PatientName: "C", ChiefComplaint: "rash", Status: "pending", Priority: "low"})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	router.ServeHTTP(w, req)

	var tasks []model.Task
	json.Unmarshal(w.Body.Bytes(), &tasks)
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}

	// Filter by status
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/tasks?status=pending", nil)
	router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &tasks)
	if len(tasks) != 2 {
		t.Errorf("expected 2 pending tasks, got %d", len(tasks))
	}
}

func TestToggleStatus(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	task := model.Task{PatientName: "D", ChiefComplaint: "pain", Status: "pending", Priority: "normal"}
	db.Create(&task)

	// pending -> in_progress
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/tasks/%d/status", task.ID), nil)
	router.ServeHTTP(w, req)

	var updated model.Task
	json.Unmarshal(w.Body.Bytes(), &updated)
	if updated.Status != "in_progress" {
		t.Errorf("expected 'in_progress', got '%s'", updated.Status)
	}

	// in_progress -> completed
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/tasks/%d/status", task.ID), nil)
	router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &updated)
	if updated.Status != "completed" {
		t.Errorf("expected 'completed', got '%s'", updated.Status)
	}
}

func TestDashboard(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	db.Create(&model.Task{PatientName: "E", ChiefComplaint: "a", Status: "pending", Priority: "normal"})
	db.Create(&model.Task{PatientName: "F", ChiefComplaint: "b", Status: "pending", Priority: "urgent"})
	db.Create(&model.Task{PatientName: "G", ChiefComplaint: "c", Status: "in_progress", Priority: "normal", TriageStatus: "completed"})
	db.Create(&model.Task{PatientName: "H", ChiefComplaint: "d", Status: "completed", Priority: "low", TriageStatus: "completed", RuleTriggered: "chest_pain"})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp DashboardResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Total != 4 {
		t.Errorf("expected total 4, got %d", resp.Total)
	}
	if resp.TriageCount != 2 {
		t.Errorf("expected triage_count 2, got %d", resp.TriageCount)
	}
	if resp.RuleOverrides != 1 {
		t.Errorf("expected rule_overrides 1, got %d", resp.RuleOverrides)
	}
}
