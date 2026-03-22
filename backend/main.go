package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/triageflow/backend/config"
	"github.com/triageflow/backend/handler"
	"github.com/triageflow/backend/llm"
	"github.com/triageflow/backend/model"
	"github.com/triageflow/backend/service"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	db, err := config.InitDB(&cfg.DB)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&model.Task{}, &model.QueueEntry{}); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	// Select triage service: LLM (EINO) or Mock
	var triageService service.Triager
	if cfg.LLM.Enabled {
		svc, err := llm.NewEinoTriageService(&cfg.LLM)
		if err != nil {
			log.Fatal("failed to init LLM triage service:", err)
		}
		triageService = svc
		log.Println("Using LLM triage service (EINO)")
	} else {
		triageService = service.NewMockTriageService()
		log.Println("Using Mock triage service")
	}
	ruleEngine := service.NewRuleEngine()
	queueService := service.NewQueueService()

	r := gin.Default()
	r.Use(cors.Default())

	taskHandler := &handler.TaskHandler{
		DB:            db,
		TriageService: triageService,
		RuleEngine:    ruleEngine,
		QueueService:  queueService,
	}
	dashHandler := &handler.DashboardHandler{DB: db}
	queueHandler := &handler.QueueHandler{DB: db, QueueService: queueService}

	api := r.Group("/api")
	{
		api.POST("/tasks", taskHandler.CreateTask)
		api.GET("/tasks", taskHandler.ListTasks)
		api.GET("/tasks/:id", taskHandler.GetTask)
		api.PATCH("/tasks/:id/status", taskHandler.ToggleStatus)
		api.GET("/dashboard", dashHandler.GetDashboard)

		api.GET("/queue", queueHandler.ListQueue)
		api.GET("/queue/:taskId/position", queueHandler.GetPosition)
		api.PATCH("/queue/:taskId/call", queueHandler.CallPatient)
		api.PATCH("/queue/:taskId/complete", queueHandler.CompletePatient)
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
