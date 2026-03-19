package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func worker(id int) {
	fmt.Printf("Worker %d starting\n", id)
	// Simulate some quick work
	time.Sleep(time.Millisecond)
	fmt.Printf("Worker %d done\n", id)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <number_of_goroutines>")
		return
	}

	// Parse the number of goroutines from the command line.
	n, err := strconv.Atoi(os.Args[1])
	if err != nil || n <= 0 {
		fmt.Println("Please provide the number of goroutines.")
		return
	}
	fmt.Printf("Starting %d goroutines...\n", n)

	// Launch goroutines without any synchronization.
	for i := 1; i <= n; i++ {
		go worker(i)
	}

	fmt.Println("Main function is exiting without waiting!")
}
