package postgresql

import (
	"database/sql"
	"exercise2/models"
	"fmt"
	"log"
)

func InsertAdvancedUsersWithTransaction(db *sql.DB, users []models.AdvancedUser) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}

	stmt, err := tx.Prepare("INSERT INTO advanced_users (name, age) VALUES ($1, $2)")
	if err != nil {
		tx.Rollback() // Rollback if there is an error preparing the statement
		return fmt.Errorf("error preparing insert statement: %v", err)
	}
	defer stmt.Close()

	for _, user := range users {
		_, err := stmt.Exec(user.Name, user.Age)
		if err != nil {
			tx.Rollback() // Rollback the transaction on error
			return fmt.Errorf("error inserting user: %v", err)
		}
	}

	err = tx.Commit() // Commit the transaction if all inserts are successful
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	fmt.Println("All users inserted successfully!")
	return nil
}

func CreateUsersTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS advanced_users (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		age INT NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	} else {
		fmt.Println("Users table created successfully!")
	}
}

func DropUsersTable(db *sql.DB) error {
	query := "DROP TABLE IF EXISTS advanced_users"

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error dropping users table: %v", err)
	}

	fmt.Println("Users table dropped successfully!")
	return nil
}

func GetAdvancedUsers(db *sql.DB, minAge int, limit int, offset int) ([]models.AdvancedUser, error) {
	query := `
	SELECT id, name, age FROM advanced_users
	WHERE age >= $1
	ORDER BY id
	LIMIT $2 OFFSET $3`

	rows, err := db.Query(query, minAge, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %v", err)
	}
	defer rows.Close()

	var users []models.AdvancedUser
	for rows.Next() {
		var user models.AdvancedUser
		if err := rows.Scan(&user.Id, &user.Name, &user.Age); err != nil {
			return nil, fmt.Errorf("error scanning user: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return users, nil
}

func UpdateAdvancedUser(db *sql.DB, id int, newName string, newAge int) error {
	query := `
	UPDATE advanced_users
	SET name = $1, age = $2
	WHERE id = $3`

	result, err := db.Exec(query, newName, newAge, id)
	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d", id)
	}

	fmt.Printf("User with ID %d updated successfully!\n", id)
	return nil
}

func DeleteAdvancedUser(db *sql.DB, id int) error {
	query := `
	DELETE FROM advanced_users WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d", id)
	}

	fmt.Printf("User with ID %d deleted successfully!\n", id)
	return nil
}
