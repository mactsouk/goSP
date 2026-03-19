package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// Record represents one JSON object in the stream
type Record struct {
	ID        int    `json:"id"`
	Event     string `json:"event"`
	Timestamp string `json:"timestamp"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:", os.Args[0], "<json-file> <sqlite-db>")
		os.Exit(1)
	}

	jsonFile := os.Args[1]
	dbFile := os.Args[2]

	// Open the JSON file
	f, err := os.Open(jsonFile)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Connect to SQLite3 database (creates it if it doesn't exist)
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY,
		event TEXT NOT NULL,
		timestamp TEXT NOT NULL
	)`)
	if err != nil {
		fmt.Println("Error creating table:", err)
		os.Exit(1)
	}

	// Prepare statement for inserting records
	stmt, err := db.Prepare(`INSERT OR IGNORE INTO events(id, event, timestamp) VALUES (?, ?, ?)`)
	if err != nil {
		fmt.Println("Error preparing insert statement:", err)
		os.Exit(1)
	}
	defer stmt.Close()

	// Create JSON decoder for streaming
	dec := json.NewDecoder(f)
	for {
		var rec Record
		if err := dec.Decode(&rec); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Skipping invalid JSON:", err)
			continue
		}

		// Insert record into database
		_, err := stmt.Exec(rec.ID, rec.Event, rec.Timestamp)
		if err != nil {
			fmt.Println("Error inserting record:", err)
		} else {
			fmt.Printf("Inserted: ID=%d | Event=%s | Time=%s\n",
				rec.ID, rec.Event, rec.Timestamp)
		}
	}
}
