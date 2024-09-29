package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func advancedTest3() {
	r := mux.NewRouter()

	// Routes for direct SQL
	r.HandleFunc("/sql/users", GetUsersSQL).Methods("GET")
	r.HandleFunc("/sql/users", CreateUserSQL).Methods("POST")
	r.HandleFunc("/sql/users/{id}", UpdateUserSQL).Methods("PUT")
	r.HandleFunc("/sql/users/{id}", DeleteUserSQL).Methods("DELETE")

	// Routes for GORM
	r.HandleFunc("/gorm/users", GetUsersGORM).Methods("GET")
	r.HandleFunc("/gorm/users", CreateUserGORM).Methods("POST")
	r.HandleFunc("/gorm/users/{id}", UpdateUserGORM).Methods("PUT")
	r.HandleFunc("/gorm/users/{id}", DeleteUserGORM).Methods("DELETE")

	// Start the server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
