package main

import (
	"encoding/json"
	"exercise2/models"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetUsersSQL(w http.ResponseWriter, r *http.Request) {
	ageFilter := r.URL.Query().Get("age")
	sortBy := r.URL.Query().Get("sortBy")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	query := "SELECT id, name, age FROM rest_users"
	var args []interface{}
	if ageFilter != "" {
		query += " WHERE age = $1"
		args = append(args, ageFilter)
	}
	if sortBy == "name" {
		query += " ORDER BY name"
	}

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	rows, err := sqlDB.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

func CreateUserSQL(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if user name is unique
	var exists int
	err = sqlDB.QueryRow("SELECT COUNT(*) FROM rest_users WHERE name = $1", user.Name).Scan(&exists)
	if err != nil {
		http.Error(w, "Error checking for user existence", http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		http.Error(w, "User with the same name already exists", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO rest_users (name, age) VALUES ($1, $2) RETURNING id`
	err = sqlDB.QueryRow(query, user.Name, user.Age).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func UpdateUserSQL(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if user exists
	var exists int
	err = sqlDB.QueryRow("SELECT COUNT(*) FROM rest_users WHERE id = $1", id).Scan(&exists)
	if err != nil {
		http.Error(w, "Error checking for user existence", http.StatusInternalServerError)
		return
	}
	if exists == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if name is unique
	err = sqlDB.QueryRow("SELECT COUNT(*) FROM rest_users WHERE name = $1 AND id != $2", user.Name, id).Scan(&exists)
	if err != nil {
		http.Error(w, "Error checking for user existence", http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		http.Error(w, "User with the same name already exists", http.StatusBadRequest)
		return
	}

	query := `UPDATE rest_users SET name = $1, age = $2 WHERE id = $3`
	_, err = sqlDB.Exec(query, user.Name, user.Age, id)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// Delete User (Direct SQL)
func DeleteUserSQL(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	// Check if user exists
	var exists int
	err := sqlDB.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", id).Scan(&exists)
	if err != nil {
		http.Error(w, "Error checking for user existence", http.StatusInternalServerError)
		return
	}
	if exists == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	query := "DELETE FROM users WHERE id = $1"
	_, err = sqlDB.Exec(query, id)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
