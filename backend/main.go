package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/triageflow/backend/config"
	"github.com/triageflow/backend/handler"
	"github.com/triageflow/backend/model"
	"github.com/triageflow/backend/service"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&model.Task{}); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	triageService := service.NewMockTriageService()
	ruleEngine := service.NewRuleEngine()

	r := gin.Default()
	r.Use(cors.Default())

	taskHandler := &handler.TaskHandler{
		DB:            db,
		TriageService: triageService,
		RuleEngine:    ruleEngine,
	}
	dashHandler := &handler.DashboardHandler{DB: db}

	api := r.Group("/api")
	{
		api.POST("/tasks", taskHandler.CreateTask)
		api.GET("/tasks", taskHandler.ListTasks)
		api.GET("/tasks/:id", taskHandler.GetTask)
		api.PATCH("/tasks/:id/status", taskHandler.ToggleStatus)
		api.GET("/dashboard", dashHandler.GetDashboard)
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
