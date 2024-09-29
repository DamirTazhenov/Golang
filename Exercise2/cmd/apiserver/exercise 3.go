package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func test3() {
	r := mux.NewRouter()
	setupRoutes(r)

	// Start server
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
