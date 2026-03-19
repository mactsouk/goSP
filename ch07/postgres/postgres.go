package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// connectDB initializes the database connection using pgx and creates a table.
func connectDB() (*sql.DB, error) {
	connStr := "postgres://postgres:password@localhost:5432/testdb?sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
                email TEXT UNIQUE NOT NULL
	)`
	if _, err := db.Exec(createTable); err != nil {
		return nil, err
	}
	return db, nil
}

// insertUser inserts a single user record.
func insertUser(db *sql.DB, name, email string) error {
	_, err := db.Exec(`INSERT INTO users (name, email) VALUES ($1, $2)`, name, email)
	return err
}

// insertUsersPrepared uses a prepared statement for multiple inserts.
func insertUsersPrepared(db *sql.DB, users map[string]string) error {
	stmt, err := db.Prepare(`INSERT INTO users (name, email) VALUES ($1, $2)`)
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

// queryUsers retrieves and prints all users.
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

// updateUser updates a user’s email by name.
func updateUser(db *sql.DB, name, newEmail string) error {
	_, err := db.Exec(`UPDATE users SET email = $1 WHERE name = $2`, newEmail, name)
	return err
}

// deleteUser removes a user by name.
func deleteUser(db *sql.DB, name string) error {
	_, err := db.Exec(`DELETE FROM users WHERE name = $1`, name)
	return err
}

func main() {
	db, err := connectDB()
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer db.Close()

	fmt.Println("Inserting users...")
	insertUser(db, "Alice", "alice@example.com")
	insertUser(db, "Bob", "bob@example.com")

	fmt.Println("\nInserting multiple users using prepared statements...")
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
