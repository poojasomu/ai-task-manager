package config

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/poojasomu/ai-task-manager/models"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate User and Task models
	err = DB.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("✅ Database connected and migrated successfully!")
}
