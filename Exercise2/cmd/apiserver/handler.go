package main

import (
	"encoding/json"
	"exercise2/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Get Users (via GORM)
func getUsersGORM(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	result := db.Find(&users)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// Create User (via GORM)
func createUserGORM(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	result := db.Create(&user)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// Update User (via GORM)
func updateUserGORM(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	result := db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	db.Save(&user)
	json.NewEncoder(w).Encode(user)
}

// Delete User (via GORM)
func deleteUserGORM(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	result := db.Delete(&models.User{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Get Users (via SQL)
func getUsersSQL(w http.ResponseWriter, r *http.Request) {
	rows, err := sqlDB.Query("SELECT id, name, age FROM users")
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			http.Error(w, "Failed to scan result", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

// Create User (via SQL)
func createUserSQL(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id`
	err = sqlDB.QueryRow(query, user.Name, user.Age).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

// Routes setup
func setupRoutes(r *mux.Router) {
	// Routes for GORM
	r.HandleFunc("/gorm/users", getUsersGORM).Methods("GET")
	r.HandleFunc("/gorm/user", createUserGORM).Methods("POST")
	r.HandleFunc("/gorm/user/{id}", updateUserGORM).Methods("PUT")
	r.HandleFunc("/gorm/user/{id}", deleteUserGORM).Methods("DELETE")

	// Routes for direct SQL queries
	r.HandleFunc("/sql/users", getUsersSQL).Methods("GET")
	r.HandleFunc("/sql/user", createUserSQL).Methods("POST")
}
