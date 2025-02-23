package handlers

import (
	"net/http"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/poojasomu/ai-task-manager/config"
	"github.com/poojasomu/ai-task-manager/models"
)

// CreateTask - API to create a new task
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := config.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// Broadcast update to WebSockets
	taskJSON, _ := json.Marshal(task)
	broadcast <- string(taskJSON)

	c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully", "task": task})
}

// GetTasks - API to fetch all tasks
func GetTasks(c *gin.Context) {
	var tasks []models.Task
	config.DB.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

// GetTask - API to fetch a single task by ID
func GetTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask - API to update a task
func UpdateTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	config.DB.Save(&task)

	// Broadcast update to WebSockets
	taskJSON, _ := json.Marshal(task)
	broadcast <- string(taskJSON)

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task": task})
}

// DeleteTask - API to delete a task
func DeleteTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")

	if err := config.DB.First(&task, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	config.DB.Delete(&task)

	// Broadcast delete event
	broadcast <- `{"message": "Task deleted", "task_id": ` + id + `}`

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
