package models

import "time"

type Team struct {
	ID          uint         `gorm:"primaryKey"`
	Name        string       `json:"name" validate:"required,min=3,max=100"`
	Description string       `json:"description" validate:"max=500"`
	CreatorID   uint         `json:"creator_id"` // ID создателя команды (менеджер)
	Members     []TeamMember `gorm:"foreignKey:TeamID"`
	Tasks       []Task       `gorm:"foreignKey:TeamID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TeamMember struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   uint   `json:"user_id"`
	TeamID   uint   `json:"team_id"`
	Role     string `json:"role" validate:"required,oneof=manager employee client"` // Роли в команде
	JoinedAt time.Time
}
