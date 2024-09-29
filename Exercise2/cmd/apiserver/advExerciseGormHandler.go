package main

import (
	"encoding/json"
	"exercise2/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Get Users (GORM)
func GetUsersGORM(w http.ResponseWriter, r *http.Request) {
	ageFilter := r.URL.Query().Get("age")
	sortBy := r.URL.Query().Get("sortBy")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var users []models.RestUser
	query := db

	if ageFilter != "" {
		query = query.Where("age = ?", ageFilter)
	}
	if sortBy == "name" {
		query = query.Order("name asc")
	}

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// Create RestUser (GORM)
func CreateUserGORM(w http.ResponseWriter, r *http.Request) {
	var user models.RestUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validation: unique name
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Error inserting user: "+err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// Update RestUser (GORM)
func UpdateUserGORM(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var user models.RestUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var existingUser models.RestUser
	if err := db.First(&existingUser, id).Error; err != nil {
		http.Error(w, "RestUser not found", http.StatusNotFound)
		return
	}

	// Validation: unique name
	if err := db.Model(&existingUser).Updates(user).Error; err != nil {
		http.Error(w, "Error updating user: "+err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(existingUser)
}

// Delete RestUser (GORM)
func DeleteUserGORM(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := db.Delete(&models.RestUser{}, id).Error; err != nil {
		http.Error(w, "Error deleting user: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
