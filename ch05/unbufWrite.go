package main

import (
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

	fmt.Printf("Writing to file '%s' using unbuffered I/O:\n", filename)
	fmt.Println("Writing byte by byte...")

	// Sample data to write
	data := "The quick brown fox jumps over the lazy dog. "
	data += "This is a demonstration of unbuffered writing in Go. "
	data += "Each byte will result in a separate system call to the operating system.\n"

	// Repeat the data to make it more substantial
	fullData := ""
	for i := 0; i < 50; i++ {
		fullData += fmt.Sprintf("Line %d: %s", i+1, data)
	}

	start := time.Now()
	bytesWritten := 0
	systemCalls := 0

	// Write the file byte by byte (unbuffered)
	for i := 0; i < len(fullData); i++ {
		buffer := []byte{fullData[i]} // Single byte buffer
		n, err := file.Write(buffer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing: %v\n", err)
			os.Exit(1)
		}
		bytesWritten += n
		systemCalls++

		// Show progress for first 50 bytes
		if i < 50 {
			if fullData[i] == '\n' {
				fmt.Print("\\n")
			} else if fullData[i] == '\t' {
				fmt.Print("\\t")
			} else if fullData[i] < 32 || fullData[i] > 126 {
				fmt.Print(".")
			} else {
				fmt.Printf("%c", fullData[i])
			}
		}
	}

	// Force data to be written to disk
	file.Sync()
	duration := time.Since(start)

	if len(fullData) > 50 {
		fmt.Print("...")
	}
	fmt.Printf("\n\nTotal bytes written: %d\n", bytesWritten)
	fmt.Printf("System calls made: %d\n", systemCalls)
	fmt.Printf("Time taken: %v\n", duration)
	fmt.Printf("Each byte required a separate write() system call.\n")
	fmt.Printf("File '%s' created successfully.\n", filename)
}
