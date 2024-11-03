package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User struct with GORM model
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password,omitempty"`
	Items    []Item `json:"items" gorm:"foreignKey:UserID"`
	RoleID   uint   `json:"role_id"` // Внешний ключ на Role
	Role     Role   `gorm:"foreignKey:RoleID"`
}

type Role struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"` // Например, "admin", "manager", "user"
}

// HashPassword hashes a user's password before storing it
func (user *User) HashPassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword compares the hashed password with the provided one
func (user *User) CheckPassword(providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
}
