package main

import (
	"bufio"
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

	fmt.Printf("Reading file '%s' using buffered I/O:\n", filename)
	fmt.Println("Using bufio.Reader with internal buffer...")

	// Create a buffered reader with default buffer size (4KB)
	reader := bufio.NewReader(file)
	byteCount := 0

	// Read the file using buffered I/O
	for {
		// ReadByte uses the internal buffer and only makes system calls
		// when the buffer needs to be refilled
		b, err := reader.ReadByte()
		if err != nil {
			break
		}

		byteCount++

		// Print first 50 bytes to show the reading process
		if byteCount <= 50 {
			if b == '\n' {
				fmt.Printf("\\n")
			} else if b == '\t' {
				fmt.Printf("\\t")
			} else if b < 32 || b > 126 {
				fmt.Printf(".")
			} else {
				fmt.Printf("%c", b)
			}
		}
	}

	if byteCount > 50 {
		fmt.Print("...")
	}
	fmt.Printf("\n\nTotal bytes read: %d\n", byteCount)
	fmt.Printf("Buffer size used: %d bytes\n", reader.Size())
}
