package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Record defines the structure of each JSON object in the stream
type Record struct {
	ID    int    `json:"id"`
	Event string `json:"event"`
	Time  string `json:"timestamp"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<json-stream-file>")
		os.Exit(1)
	}

	// Open the JSON stream file
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Create a new JSON decoder for the stream
	dec := json.NewDecoder(f)

	// Read each JSON object in the stream
	for {
		var rec Record
		if err := dec.Decode(&rec); err != nil {
			// When the stream ends, json.Decoder returns io.EOF
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Skipping invalid JSON:", err)
			continue
		}

		// Process each decoded record
		fmt.Printf("ID: %d | Event: %s | Time: %s\n", rec.ID, rec.Event, rec.Time)
	}
}
