package main

import (
	"exercise2/internal/config"
	"exercise2/models"
	"exercise2/pkg/postgresql"
	"fmt"
	"log"
)

func printUsers(users []models.AdvancedUser) {
	for _, user := range users {
		fmt.Printf("Name: %s, Age: %d\n", user.Name, user.Age)
	}
}

func advancedTest(configData *config.Config) {
	// Initialize DB connection
	db := postgresql.InitAdvancedDB(configData)
	defer db.Close()

	postgresql.DropUsersTable(db)
	// Create users table
	postgresql.CreateUsersTable(db)

	// Insert users with transaction
	users := []models.AdvancedUser{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}
	err := postgresql.InsertAdvancedUsersWithTransaction(db, users)
	if err != nil {
		log.Fatalf("Error inserting users: %v", err)
	}

	// Query users with pagination and filtering
	queriedUsers, err := postgresql.GetAdvancedUsers(db, 25, 2, 0)
	if err != nil {
		log.Fatalf("Error querying users: %v", err)
	}
	printUsers(queriedUsers)

	// Update user
	err = postgresql.UpdateAdvancedUser(db, 1, "Alice Updated", 32)
	if err != nil {
		log.Fatalf("Error updating user: %v", err)
	}

	// Delete user
	err = postgresql.DeleteAdvancedUser(db, 2)
	if err != nil {
		log.Fatalf("Error deleting user: %v", err)
	}
}
