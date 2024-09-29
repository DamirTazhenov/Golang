package postgresql

import (
	"database/sql"
	"exercise2/internal/config"
	"fmt"
	"log"
	"time"
)

func InitAdvancedDB(configData *config.Config) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		configData.Host, configData.Port, configData.User, configData.Password, configData.DbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Set connection pooling settings
	db.SetMaxOpenConns(25)                 // Maximum number of open connections
	db.SetMaxIdleConns(25)                 // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Max lifetime of a connection

	// Ping to ensure the connection is established
	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}
	return db
}
