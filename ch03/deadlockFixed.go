package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== Fixed Deadlock Example ===")
	ch := make(chan int)

	// Worker goroutine sends a value and exits.
	go func() {
		fmt.Println("Goroutine: sending value...")
		ch <- 42
		fmt.Println("Goroutine: done")
	}()

	// Main goroutine receives the value.
	time.Sleep(time.Second)
	fmt.Println("Main: receiving value...")
	val := <-ch
	fmt.Println("Main: received", val)

	fmt.Println("End of program")
}
