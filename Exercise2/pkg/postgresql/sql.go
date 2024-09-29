package postgresql

import (
	"database/sql"
	"fmt"
	"log"
)

func CreateTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		age INT
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	} else {
		fmt.Println("Table created successfully!")
	}
}

func InsertUser(db *sql.DB, name string, age int) {
	query := `
	INSERT INTO users (name, age)
	VALUES ($1, $2)
	RETURNING id;`

	var id int
	err := db.QueryRow(query, name, age).Scan(&id)
	if err != nil {
		log.Fatalf("Error inserting user: %v", err)
	}
	fmt.Printf("Inserted user with ID: %d\n", id)
}

func GetUsers(db *sql.DB) {
	query := `SELECT id, name, age FROM users;`

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Error retrieving users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var age int

		err := rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func closeDB(sqlDb *sql.DB) {
	sqlDb.Close()
}
