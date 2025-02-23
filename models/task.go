package models

import "gorm.io/gorm"

type Task struct {
    gorm.Model
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      string `json:"status" gorm:"default:'Pending'"`
    AssignedTo  string `json:"assigned_to"`
}
