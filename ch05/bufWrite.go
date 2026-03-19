package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	// Check if filename argument is provided
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <filename>\n", os.Args[0])
		os.Exit(1)
	}

	filename := os.Args[1]

	// Create or truncate the file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Printf("Writing to file '%s' using buffered I/O:\n", filename)
	fmt.Println("Using bufio.Writer with internal buffer...")

	// Create a buffered writer with default buffer size (4KB)
	writer := bufio.NewWriter(file)
	defer writer.Flush() // Ensure all data is written

	// Sample data to write
	data := "The quick brown fox jumps over the lazy dog. "
	data += "This is a demonstration of buffered writing in Go. "
	data += "Bytes are accumulated in a buffer before being written to the OS.\n"

	// Repeat the data to make it more substantial
	fullData := ""
	for i := 0; i < 50; i++ {
		fullData += fmt.Sprintf("Line %d: %s", i+1, data)
	}

	start := time.Now()
	bytesWritten := 0

	// Write the file byte by byte using buffered I/O
	for i := 0; i < len(fullData); i++ {
		err := writer.WriteByte(fullData[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing: %v\n", err)
			os.Exit(1)
		}
		bytesWritten++
	}

	// Flush the buffer to ensure all data is written
	writer.Flush()
	file.Sync() // Force data to be written to disk
	duration := time.Since(start)

	if len(fullData) > 50 {
		fmt.Print("...")
	}
	fmt.Printf("\n\nTotal bytes written: %d\n", bytesWritten)
	fmt.Printf("Buffer size used: %d bytes\n", writer.Size())
	fmt.Printf("Time taken: %v\n", duration)
	fmt.Println("System calls were minimized by accumulating data in buffer.")

	// Alternative demonstration: writing strings and lines
	fmt.Println("\n--- Alternative: String and line writing ---")

	// Create a new buffered writer for demonstration
	file.Seek(0, 0)  // Reset to beginning (though we're appending)
	file.Truncate(0) // Clear the file
	lineWriter := bufio.NewWriter(file)
	defer lineWriter.Flush()

	// Write some strings efficiently
	lines := []string{
		"This is line 1 written as a complete string",
		"This is line 2 with efficient string writing",
		"Line 3 demonstrates buffered string operations",
		"Line 4 shows how bufio handles larger chunks",
		"Final line written efficiently to buffer",
	}

	linesWritten := 0
	for _, line := range lines {
		n, err := lineWriter.WriteString(line + "\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing line: %v\n", err)
			break
		}
		linesWritten++
		fmt.Printf("Wrote line %d (%d bytes)\n", linesWritten, n)
	}

	lineWriter.Flush()
	fmt.Printf("Total lines written: %d\n", linesWritten)
	fmt.Printf("File '%s' created successfully.\n", filename)
}
