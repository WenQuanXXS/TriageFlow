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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	host := getEnvOrDefault("DB_HOST", "127.0.0.1")
	port := getEnvOrDefault("DB_PORT", "3306")
	user := getEnvOrDefault("DB_USER", "root")
	pass := getEnvOrDefault("DB_PASS", "root123")
	name := getEnvOrDefault("DB_NAME", "triageflow_test")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("skipping: cannot connect to test database: %v", err)
	}

	db.Migrator().DropTable(&model.Task{})
	db.AutoMigrate(&model.Task{})

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

	taskHandler := &TaskHandler{DB: db}
	dashHandler := &DashboardHandler{DB: db}

	api := r.Group("/api")
	{
		api.POST("/tasks", taskHandler.CreateTask)
		api.GET("/tasks", taskHandler.ListTasks)
		api.PATCH("/tasks/:id/status", taskHandler.ToggleStatus)
		api.GET("/dashboard", dashHandler.GetDashboard)
	}

	return r
}

func TestCreateTask(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	body := `{"patient_name":"John Doe","chief_complaint":"headache","priority":"urgent","department":"Neurology"}`
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var task model.Task
	json.Unmarshal(w.Body.Bytes(), &task)

	if task.PatientName != "John Doe" {
		t.Errorf("expected patient_name 'John Doe', got '%s'", task.PatientName)
	}
	if task.Status != "pending" {
		t.Errorf("expected status 'pending', got '%s'", task.Status)
	}
	if task.Priority != "urgent" {
		t.Errorf("expected priority 'urgent', got '%s'", task.Priority)
	}
}

func TestListTasks(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	// Seed data
	db.Create(&model.Task{PatientName: "A", ChiefComplaint: "cough", Status: "pending", Priority: "normal"})
	db.Create(&model.Task{PatientName: "B", ChiefComplaint: "fever", Status: "in_progress", Priority: "urgent"})
	db.Create(&model.Task{PatientName: "C", ChiefComplaint: "rash", Status: "pending", Priority: "low"})

	// List all
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

	// Filter by priority
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/tasks?priority=urgent", nil)
	router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &tasks)
	if len(tasks) != 1 {
		t.Errorf("expected 1 urgent task, got %d", len(tasks))
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

	// completed -> pending
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/tasks/%d/status", task.ID), nil)
	router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &updated)
	if updated.Status != "pending" {
		t.Errorf("expected 'pending', got '%s'", updated.Status)
	}
}

func TestDashboard(t *testing.T) {
	db := setupTestDB(t)
	router := setupRouter(db)

	// Seed mixed data
	db.Create(&model.Task{PatientName: "E", ChiefComplaint: "a", Status: "pending", Priority: "normal"})
	db.Create(&model.Task{PatientName: "F", ChiefComplaint: "b", Status: "pending", Priority: "urgent"})
	db.Create(&model.Task{PatientName: "G", ChiefComplaint: "c", Status: "in_progress", Priority: "normal"})
	db.Create(&model.Task{PatientName: "H", ChiefComplaint: "d", Status: "completed", Priority: "low"})

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

	statusMap := map[string]int64{}
	for _, s := range resp.ByStatus {
		statusMap[s.Status] = s.Count
	}
	if statusMap["pending"] != 2 {
		t.Errorf("expected 2 pending, got %d", statusMap["pending"])
	}
	if statusMap["in_progress"] != 1 {
		t.Errorf("expected 1 in_progress, got %d", statusMap["in_progress"])
	}
	if statusMap["completed"] != 1 {
		t.Errorf("expected 1 completed, got %d", statusMap["completed"])
	}
}
