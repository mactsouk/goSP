package main

import (
	"errors"
	"fmt"
	"os"
)

// Custom error type for division errors
type DivideError struct {
	A, B float64
}

func (e *DivideError) Error() string {
	return fmt.Sprintf("cannot divide %.2f by %.2f", e.A, e.B)
}

// divide performs division and demonstrates basic error handling
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, &DivideError{A: a, B: b}
	}
	return a / b, nil
}

// riskyFileOperation simulates opening a file and panicking if it fails
func riskyFileOperation(filename string) {
	defer func() {
		// Recover from panic if it occurs
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()

	// Try to open a file that may not exist
	f, err := os.Open(filename)
	if err != nil {
		// Wrapping the error with additional context
		panic(fmt.Errorf("critical failure opening %q: %w", filename, err))
	}
	defer f.Close()

	fmt.Println("Successfully opened file:", filename)
}

func main() {
	fmt.Println("=== Basic Error Handling ===")
	result, err := divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)

		// Using errors.As to check for custom error type
		var divErr *DivideError
		if errors.As(err, &divErr) {
			fmt.Printf("Info: A=%.2f, B=%.2f\n", divErr.A, divErr.B)
		}
	} else {
		fmt.Println("Result:", result)
	}

	fmt.Println("\n=== Using errors.Is ===")
	// Simulate file handling and check against known sentinel errors
	_, err = os.Open("/tmp/nonexistent.txt")
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("The file does not exist.")
	}

	fmt.Println("\n=== Panic, Recover, and Defer ===")
	// Attempt a risky file operation
	riskyFileOperation("/tmp/nonexistent.txt")
	fmt.Println("Program continues normally after recover.")
}
