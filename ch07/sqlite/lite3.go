package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// connectDB opens or creates a SQLite database and ensures the users table exists.
func connectDB(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT
	)`); err != nil {
		return nil, err
	}
	return db, nil
}

// insertUser adds a new record using Exec.
func insertUser(db *sql.DB, name, email string) error {
	_, err := db.Exec(`INSERT INTO users (name, email) VALUES (?, ?)`, name, email)
	return err
}

// insertUsersPrepared inserts multiple records using a prepared statement.
func insertUsersPrepared(db *sql.DB, users map[string]string) error {
	stmt, err := db.Prepare(`INSERT INTO users (name, email) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for name, email := range users {
		if _, err := stmt.Exec(name, email); err != nil {
			return err
		}
	}
	return nil
}

// queryUsers retrieves and prints all user records.
func queryUsers(db *sql.DB) error {
	rows, err := db.Query(`SELECT id, name, email FROM users ORDER BY id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("Current users:")
	for rows.Next() {
		var id int
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			return err
		}
		fmt.Printf("%d | %s | %s\n", id, name, email)
	}
	return rows.Err()
}

// updateUser modifies a user’s email by name.
func updateUser(db *sql.DB, name, newEmail string) error {
	_, err := db.Exec(`UPDATE users SET email = ? WHERE name = ?`, newEmail, name)
	return err
}

// deleteUser removes a user by name.
func deleteUser(db *sql.DB, name string) error {
	_, err := db.Exec(`DELETE FROM users WHERE name = ?`, name)
	return err
}

func main() {
	db, err := connectDB("example_prepared.db")
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer db.Close()

	fmt.Println("Inserting a few users individually...")
	insertUser(db, "Alice", "alice@example.com")
	insertUser(db, "Bob", "bob@example.com")

	fmt.Println("\nInserting multiple users using a prepared statement...")
	users := map[string]string{
		"Charlie": "charlie@example.com",
		"Dave":    "dave@example.com",
		"Eve":     "eve@example.com",
	}
	if err := insertUsersPrepared(db, users); err != nil {
		log.Fatal("Prepared insert error:", err)
	}

	queryUsers(db)

	fmt.Println("\nUpdating Bob’s email...")
	updateUser(db, "Bob", "bob@newdomain.com")

	queryUsers(db)

	fmt.Println("\nDeleting Alice...")
	deleteUser(db, "Alice")

	queryUsers(db)
}
