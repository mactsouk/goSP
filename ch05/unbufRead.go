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

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Printf("Reading file '%s' using unbuffered I/O:\n", filename)
	fmt.Println("Reading byte by byte...")

	// Read the file unbuffered
	buffer := make([]byte, 16) // 16-byte buffer
	byteCount := 0

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			byteCount++
			// Print first 16x bytes to show the reading process
			if byteCount <= 50 {
				if buffer[0] == '\n' {
					fmt.Printf("\\n")
				} else if buffer[0] == '\t' {
					fmt.Printf("\\t")
				} else if buffer[0] < 32 || buffer[0] > 126 {
					fmt.Printf(".")
				} else {
					fmt.Printf("%c", buffer[0])
				}
			}
		}
		if err != nil {
			break
		}
	}

	if byteCount > 50 {
		fmt.Print("...")
	}
	fmt.Printf("\n\nTotal bytes read: %d\n", byteCount)
}
