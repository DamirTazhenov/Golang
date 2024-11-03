package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"not null"`
	Price  float64
	UserID uint
}

var NextID uint = 1
