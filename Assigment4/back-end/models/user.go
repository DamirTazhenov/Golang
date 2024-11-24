package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User struct with GORM model
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password,omitempty"`
	Tasks    []Task `json:"tasks" gorm:"foreignKey:UserID"`
	Mode     string `json:"mode" validate:"required,oneof=individual team"`
}

// HashPassword hashes a user's password before storing it
func (user *User) HashPassword(password string) error {
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
