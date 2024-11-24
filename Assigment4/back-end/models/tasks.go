package models

import (
	"gorm.io/gorm"
)

// Task struct with GORM model
type Task struct {
	gorm.Model
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required,min=10,max=500"`
	Completed   bool   `json:"completed"`
	UserID      uint   `json:"user_id"`
	TeamID      uint   `json:"team_id"`
}
