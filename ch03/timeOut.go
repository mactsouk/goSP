package main

import (
	"fmt"
	"time"
)

// Simulates a long-running task (5 seconds).
func longTask(done chan<- string) {
	time.Sleep(5 * time.Second)
	done <- "Long task finished"
}

// Simulates a short-running task (1 second).
func shortTask(done chan<- string) {
	time.Sleep(1 * time.Second)
	done <- "Short task finished"
}

func main() {
	longDone := make(chan string)
	shortDone := make(chan string)

	// Start both tasks.
	go longTask(longDone)
	go shortTask(shortDone)

	// Handle the short task with a 2-second timeout.
	select {
	case result := <-shortDone:
		fmt.Println("Short:", result)
	case <-time.After(2 * time.Second):
		fmt.Println("Short: Timeout! Took too long.")
	}

	// Handle the long task with a 2-second timeout.
	select {
	case result := <-longDone:
		fmt.Println("Long:", result)
	case <-time.After(2 * time.Second):
		fmt.Println("Long: Timeout! Took too long.")
	}
}
