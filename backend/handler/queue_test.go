package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/triageflow/backend/model"
	"github.com/triageflow/backend/service"
	"gorm.io/gorm"
)

func setupQueueTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := setupTestDB(t)
	db.Migrator().DropTable(&model.QueueEntry{})
	db.AutoMigrate(&model.QueueEntry{})
	return db
}

func setupQueueRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	queueSvc := service.NewQueueService()

	taskHandler := &TaskHandler{
		DB:            db,
		TriageService: service.NewMockTriageService(),
		RuleEngine:    service.NewRuleEngine(),
		QueueService:  queueSvc,
	}
	queueHandler := &QueueHandler{DB: db, QueueService: queueSvc}

	api := r.Group("/api")
	{
		api.POST("/tasks", taskHandler.CreateTask)
		api.GET("/tasks/:id", taskHandler.GetTask)
		api.GET("/queue", queueHandler.ListQueue)
		api.GET("/queue/:taskId/position", queueHandler.GetPosition)
		api.PATCH("/queue/:taskId/call", queueHandler.CallPatient)
		api.PATCH("/queue/:taskId/complete", queueHandler.CompletePatient)
	}

	return r
}

func createTestTask(t *testing.T, router *gin.Engine, name, complaint string) model.Task {
	t.Helper()
	body := fmt.Sprintf(`{"patient_name":"%s","chief_complaint":"%s"}`, name, complaint)
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var task model.Task
	json.Unmarshal(w.Body.Bytes(), &task)
	return task
}

func TestAutoEnqueue(t *testing.T) {
	db := setupQueueTestDB(t)
	router := setupQueueRouter(db)

	task := createTestTask(t, router, "患者A", "头痛三天")

	var entry model.QueueEntry
	if err := db.Where("task_id = ?", task.ID).First(&entry).Error; err != nil {
		t.Fatalf("expected queue entry for task %d, got error: %v", task.ID, err)
	}

	if entry.PatientName != "患者A" {
		t.Errorf("expected patient_name '患者A', got '%s'", entry.PatientName)
	}
	if entry.QueueStatus != "waiting" {
		t.Errorf("expected queue_status 'waiting', got '%s'", entry.QueueStatus)
	}
	if entry.Department == "" {
		t.Error("expected department to be set")
	}
	if entry.QueueNumber != 1 {
		t.Errorf("expected queue_number 1, got %d", entry.QueueNumber)
	}
}

func TestNoDuplicateEnqueue(t *testing.T) {
	db := setupQueueTestDB(t)
	router := setupQueueRouter(db)

	task := createTestTask(t, router, "患者B", "咳嗽两天")

	// Try to enqueue again via the service directly
	queueSvc := service.NewQueueService()
	err := queueSvc.Enqueue(db, &task)
	if err != nil {
		t.Fatalf("unexpected error on duplicate enqueue: %v", err)
	}

	var count int64
	db.Model(&model.QueueEntry{}).Where("task_id = ?", task.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 queue entry, got %d", count)
	}
}

func TestQueueOrdering_PriorityFirst(t *testing.T) {
	db := setupQueueTestDB(t)
	router := setupQueueRouter(db)

	// normal priority: 头痛 (headache, no risk)
	createTestTask(t, router, "Normal患者", "头痛")
	time.Sleep(10 * time.Millisecond)
	// urgent priority: 胸痛 (chest pain, triggers rule engine → urgent)
	createTestTask(t, router, "Urgent患者", "胸痛")
	time.Sleep(10 * time.Millisecond)
	// high priority: 高热 triggers rule → high
	createTestTask(t, router, "High患者", "高热")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/queue", nil)
	router.ServeHTTP(w, req)

	var entries []model.QueueEntry
	json.Unmarshal(w.Body.Bytes(), &entries)

	if len(entries) < 3 {
		t.Fatalf("expected at least 3 queue entries, got %d", len(entries))
	}

	// urgent should come first globally
	if entries[0].Priority != "urgent" {
		t.Errorf("expected first entry to be urgent, got '%s' (patient: %s)", entries[0].Priority, entries[0].PatientName)
	}
}

func TestQueueOrdering_SamePriorityByTime(t *testing.T) {
	db := setupQueueTestDB(t)
	router := setupQueueRouter(db)

	task1 := createTestTask(t, router, "First", "头痛")
	time.Sleep(50 * time.Millisecond)
	task2 := createTestTask(t, router, "Second", "头晕")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/queue?department=Neurology", nil)
	router.ServeHTTP(w, req)

	var entries []model.QueueEntry
	json.Unmarshal(w.Body.Bytes(), &entries)

	if len(entries) < 2 {
		t.Fatalf("expected at least 2 entries, got %d", len(entries))
	}

	if entries[0].TaskID != task1.ID {
		t.Errorf("expected first entry task_id=%d, got %d", task1.ID, entries[0].TaskID)
	}
	if entries[1].TaskID != task2.ID {
		t.Errorf("expected second entry task_id=%d, got %d", task2.ID, entries[1].TaskID)
	}
}

func TestCallPatient(t *testing.T) {
	db := setupQueueTestDB(t)
	router := setupQueueRouter(db)

	task := createTestTask(t, router, "叫号患者", "头痛")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/queue/%d/call", task.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var entry model.QueueEntry
	json.Unmarshal(w.Body.Bytes(), &entry)

	if entry.QueueStatus != "called" {
		t.Errorf("expected queue_status 'called', got '%s'", entry.QueueStatus)
	}
	if entry.CalledAt == nil {
		t.Error("expected called_at to be set")
	}

	var updatedTask model.Task
	db.First(&updatedTask, task.ID)
	if updatedTask.Status != "in_progress" {
		t.Errorf("expected task status 'in_progress', got '%s'", updatedTask.Status)
	}
}

func TestCallPatient_AlreadyCalled(t *testing.T) {
	db := setupQueueTestDB(t)
	router := setupQueueRouter(db)

	task := createTestTask(t, router, "重复叫号", "头痛")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/queue/%d/call", task.ID), nil)
	router.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/queue/%d/call", task.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for double call, got %d", w.Code)
	}
}

func TestCompletePatient(t *testing.T) {
	db := setupQueueTestDB(t)
	router := setupQueueRouter(db)

	task := createTestTask(t, router, "完成患者", "头痛")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/queue/%d/call", task.ID), nil)
	router.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/queue/%d/complete", task.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var entry model.QueueEntry
	json.Unmarshal(w.Body.Bytes(), &entry)

	if entry.QueueStatus != "completed" {
		t.Errorf("expected queue_status 'completed', got '%s'", entry.QueueStatus)
	}
	if entry.CompletedAt == nil {
		t.Error("expected completed_at to be set")
	}

	var updatedTask model.Task
	db.First(&updatedTask, task.ID)
	if updatedTask.Status != "completed" {
		t.Errorf("expected task status 'completed', got '%s'", updatedTask.Status)
	}
}

func TestCompletePatient_NotCalled(t *testing.T) {
	db := setupQueueTestDB(t)
	router := setupQueueRouter(db)

	task := createTestTask(t, router, "未叫号", "头痛")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/queue/%d/complete", task.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetPosition(t *testing.T) {
	db := setupQueueTestDB(t)
	router := setupQueueRouter(db)

	task1 := createTestTask(t, router, "先来1", "头痛")
	time.Sleep(50 * time.Millisecond)
	task2 := createTestTask(t, router, "先来2", "头晕")
	time.Sleep(50 * time.Millisecond)
	task3 := createTestTask(t, router, "后来3", "偏头痛")

	// Check position of task3
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/queue/%d/position", task3.ID), nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var pos service.PositionInfo
	json.Unmarshal(w.Body.Bytes(), &pos)

	if pos.Ahead != 2 {
		t.Errorf("expected 2 ahead, got %d", pos.Ahead)
	}
	if pos.QueueStatus != "waiting" {
		t.Errorf("expected 'waiting', got '%s'", pos.QueueStatus)
	}
	if len(pos.WaitingList) < 3 {
		t.Errorf("expected at least 3 in waiting_list, got %d", len(pos.WaitingList))
	}
	if len(pos.NowServing) != 0 {
		t.Errorf("expected 0 in now_serving, got %d", len(pos.NowServing))
	}

	// Call task1 → should appear in now_serving, ahead count decreases
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/queue/%d/call", task1.ID), nil)
	router.ServeHTTP(w, req)

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/queue/%d/position", task3.ID), nil)
	router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &pos)
	if pos.Ahead != 1 {
		t.Errorf("after calling task1, expected 1 ahead, got %d", pos.Ahead)
	}
	if len(pos.NowServing) != 1 {
		t.Errorf("expected 1 in now_serving, got %d", len(pos.NowServing))
	}

	// Check task2 position
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/queue/%d/position", task2.ID), nil)
	router.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &pos)
	if pos.Ahead != 0 {
		t.Errorf("after calling task1, expected 0 ahead for task2, got %d", pos.Ahead)
	}
}

func TestQueueFilterByDepartment(t *testing.T) {
	db := setupQueueTestDB(t)
	router := setupQueueRouter(db)

	createTestTask(t, router, "神经科患者", "头痛")
	createTestTask(t, router, "呼吸科患者", "咳嗽")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/queue?department=Neurology", nil)
	router.ServeHTTP(w, req)

	var entries []model.QueueEntry
	json.Unmarshal(w.Body.Bytes(), &entries)

	for _, e := range entries {
		if e.Department != "Neurology" {
			t.Errorf("expected all entries to be Neurology, got '%s'", e.Department)
		}
	}
}
