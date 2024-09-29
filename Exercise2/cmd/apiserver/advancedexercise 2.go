package main

import (
	"exercise2/pkg/postgresql"
	"gorm.io/gorm"
	"log"
)

func advancedTest2(db *gorm.DB) {
	if err := postgresql.InsertUserWithProfile(db); err != nil {
		log.Fatalf("Error inserting user with profile: %v", err)
	}

	// Query users with profiles
	_, err := postgresql.GetUsersWithProfiles(db)
	if err != nil {
		log.Fatalf("Error fetching users with profiles: %v", err)
	}

	// Update a user's profile
	err = postgresql.UpdateUserProfile(db, 1, "Senior Developer", "http://example.com/new_alice.jpg")
	if err != nil {
		log.Fatalf("Error updating user profile: %v", err)
	}

	// Delete a user and profile
	err = postgresql.DeleteUserWithProfile(db, 1)
	if err != nil {
		log.Fatalf("Error deleting user: %v", err)
	}
}
