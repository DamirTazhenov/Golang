package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string
	Age  int
}

type AdvancedUser struct {
	Id   int
	Name string
	Age  int
}

type AdvancedGormUser struct {
	ID      uint                `gorm:"primaryKey"`
	Name    string              `gorm:"not null"`
	Age     int                 `gorm:"not null"`
	Profile AdvancedGormProfile `gorm:"foreignKey:AdvancedGormUserID;constraint:OnDelete:CASCADE"`
}

type AdvancedGormProfile struct {
	ID                 uint   `gorm:"primaryKey"`
	AdvancedGormUserID uint   `gorm:"unique;not null"` // ForeignKey explicitly defined as AdvancedGormUserID
	Bio                string `gorm:"size:255"`
	ProfilePictureURL  string
}

type RestUser struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`
	Age  int    `json:"age"`
}
