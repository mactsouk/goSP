package main

import (
	"bufio"
	"encoding/json/v2"
	"fmt"
	"os"
	"strings"
)

// Record defines the structure of each JSON object
type Record struct {
	ID    int    `json:"id"`
	Event string `json:"event"`
	Time  string `json:"timestamp"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<json-file>")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var inArray bool
	var arrayBuffer strings.Builder

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		// Detect start and end of a JSON array
		if strings.HasPrefix(line, "[") {
			inArray = true
			line = line[1:] // remove [
		}
		if strings.HasSuffix(line, "]") {
			inArray = false
			line = line[:len(line)-1] // remove ]
		}

		// Remove trailing commas
		line = strings.TrimSuffix(line, ",")

		if inArray {
			arrayBuffer.WriteString(line)
			// End of an object in the array?
			if strings.HasSuffix(line, "}") {
				processRecord(arrayBuffer.String())
				arrayBuffer.Reset()
			}
		} else {
			processRecord(line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

// processRecord unmarshals a single JSON object and prints it
func processRecord(data string) {
	var rec Record
	if err := json.Unmarshal([]byte(data), &rec); err != nil {
		fmt.Println("Skipping invalid JSON:", err)
		return
	}
	fmt.Printf("ID: %d | Event: %s | Time: %s\n", rec.ID, rec.Event, rec.Time)
}
