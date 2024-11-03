package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"net/http"
	"shop/models"
	"strconv"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	result := DB.Create(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// GetUser - Retrieve a user by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var user models.User
	result := DB.Preload("Tasks").First(&user, id)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Роль не найдена"})
			c.Abort()
			return
		}

		// Проверка, есть ли роль пользователя в списке разрешенных ролей
		userRole := role.(string)
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next() // Разрешить доступ, если роль совпадает
				return
			}
		}

		// Если ни одна из ролей не совпала, отказать в доступе
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied for this role"})
		c.Abort()
	}
}
