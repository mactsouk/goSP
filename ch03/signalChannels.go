package main

import (
	"fmt"
	"time"
)

func worker(done chan struct{}) {
	fmt.Println("Worker: starting work")
	time.Sleep(2 * time.Second)
	fmt.Println("Worker: done with work")
	// Signal completion by sending an empty struct
	done <- struct{}{}
}

func main() {
	// Channel used only for signaling
	done := make(chan struct{})
	go worker(done)

	// Wait for the signal before exiting
	<-done
	fmt.Println("Main: received signal, exiting")
}
