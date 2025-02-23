package main

import (
	"github.com/gin-gonic/gin"
	"github.com/poojasomu/ai-task-manager/config"
	"github.com/poojasomu/ai-task-manager/handlers"
	"github.com/poojasomu/ai-task-manager/middleware"
)

func main() {
	r := gin.Default()

	// Initialize Database
	config.ConnectDB()

	// Start WebSocket broadcaster in a goroutine
	go handlers.StartBroadcast()

	// Home Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "AI Task Manager API is running!"})
	})

	// Public Routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// WebSocket Route
	r.GET("/ws", handlers.HandleWebSocket)

	// Protected Routes (JWT Required)
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())

	// Task Management Routes
	protected.POST("/tasks", handlers.CreateTask)
	protected.GET("/tasks", handlers.GetTasks)
	protected.GET("/tasks/:id", handlers.GetTask)
	protected.PUT("/tasks/:id", handlers.UpdateTask)
	protected.DELETE("/tasks/:id", handlers.DeleteTask)

	// Dashboard Route
	protected.GET("/dashboard", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the dashboard!"})
	})

	// Start Server
	r.Run(":8080")
}
