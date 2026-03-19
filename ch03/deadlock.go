package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Demonstrating a Deadlock ===")
	ch := make(chan int)

	// Worker goroutine sends a value, but no one is receiving.
	go func() {
		fmt.Println("Goroutine: sending value...")
		ch <- 42 // This blocks forever
		fmt.Println("Goroutine: done")
	}()

	// Main goroutine also tries to send instead of receiving.
	time.Sleep(time.Second)
	fmt.Println("Main: also sending value...")
	ch <- 99 // Deadlock occurs here

	// This line is never reached.
	fmt.Println("End of program")
}
