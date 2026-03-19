package main

import (
	"fmt"
	"time"
)

// producer sends numbers to the channel and then closes it.
func producer(ch chan int, count int) {
	for i := 1; i <= count; i++ {
		fmt.Printf("Producer: sending %d\n", i)
		ch <- i
		time.Sleep(time.Millisecond * 300)
	}
	fmt.Println("Producer: closing channel")
	close(ch)
}

// consumer receives numbers from the channel until it’s closed.
func consumer(ch chan int) {
	for val := range ch {
		fmt.Printf("Consumer: received %d\n", val)
		time.Sleep(time.Millisecond * 500)
	}
	fmt.Println("Consumer: channel closed, stopping")
}

func main() {
	fmt.Println("=== Unbuffered Channel Example ===")
	unbuffered := make(chan int)

	// Start producer and consumer goroutines.
	go producer(unbuffered, 5)
	consumer(unbuffered) // Runs in main goroutine.

	fmt.Println("\n=== Buffered Channel Example ===")
	buffered := make(chan string, 3)

	// Demonstrate buffered channel behavior.
	fmt.Println("Main: sending A")
	buffered <- "A"
	fmt.Println("Main: sending B")
	buffered <- "B"
	fmt.Println("Main: sending C")
	buffered <- "C"

	fmt.Println("Main: retrieving values...")
	fmt.Println("Received:", <-buffered)
	fmt.Println("Received:", <-buffered)
	fmt.Println("Received:", <-buffered)

	fmt.Println("*** Done!")
}
