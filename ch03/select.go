package main

import (
	"fmt"
    "math/rand"
	"time"
)

func worker(name string, ch chan<- string) {
	// Simulate variable work duration
	delay := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(delay)
	ch <- fmt.Sprintf("%s finished after %v", name, delay)
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go worker("Worker 1", ch1)
	go worker("Worker 2", ch2)

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received:", msg2)
		}
	}

	fmt.Println("All workers reported back!")
}
