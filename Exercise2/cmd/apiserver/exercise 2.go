package main

import (
	"exercise2/models"
	"fmt"
	"gorm.io/gorm"
	"log"
)

func insertUser(db *gorm.DB, name string, age int) {
	user := models.User{Name: name, Age: age}
	result := db.Create(&user)

	if result.Error != nil {
		log.Fatalf("Failed to insert user: %v", result.Error)
	}

	fmt.Printf("Inserted user with ID: %d\n", user.ID)
}

func getUsers(db *gorm.DB) {
	var users []models.User
	result := db.Find(&users)

	if result.Error != nil {
		log.Fatalf("Failed to retrieve users: %v", result.Error)
	}

	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}

func test2(db *gorm.DB) {
	insertUser(db, "Alice", 30)
	insertUser(db, "Bob", 25)
	insertUser(db, "Charlie", 35)

	fmt.Println("Users in the database:")
	getUsers(db)
}
