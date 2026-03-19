package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// Record defines the structure of the JSON input and XML output
type Record struct {
	ID    int    `json:"id" xml:"id"`
	Event string `json:"event" xml:"event"`
	Time  string `json:"timestamp" xml:"timestamp"`
}

// Records is used to wrap multiple records for XML
type Records struct {
	XMLName xml.Name `xml:"records"`
	Items   []Record `xml:"record"`
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

	var records []Record
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var rec Record
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			fmt.Println("Skipping invalid JSON:", err)
			continue
		}

		records = append(records, rec)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	// Wrap the records and marshal to XML
	output := Records{Items: records}
	xmlData, err := xml.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Println("Error generating XML:", err)
		os.Exit(1)
	}

	fmt.Println(xml.Header + string(xmlData))
}
