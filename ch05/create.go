package main

import (
	"fmt"
	"os"
)

func main() {
	// Check if filename argument is provided
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}
	filename := os.Args[1]

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file", filename, err)
		os.Exit(1)
	}

	err = file.Close()
	if err != nil {
		fmt.Println("Error creating file", filename, err)
		os.Exit(1)
	}

	// Check if the file was a new creation or an override
	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking file info: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created empty file: %s\n", filename)
	fmt.Printf("File size: %d bytes\n", fileInfo.Size())
	fmt.Printf("File permissions: %v\n", fileInfo.Mode())
}
