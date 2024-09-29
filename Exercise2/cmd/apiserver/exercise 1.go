package main

import (
	"database/sql"
	"exercise2/internal/config"
	"exercise2/pkg/postgresql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func test(configData *config.Config) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		configData.Host, configData.Port, configData.User, configData.Password, configData.DbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	postgresql.CreateTable(db)

	postgresql.InsertUser(db, "Alice", 30)
	postgresql.InsertUser(db, "Bob", 25)
	postgresql.InsertUser(db, "Charlie", 35)

	fmt.Println("Users in the database:")
	postgresql.GetUsers(db)
}
