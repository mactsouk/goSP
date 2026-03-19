package main

import (
	"fmt"
	"time"
)

// sender generates numbers and sends them over the channel.
func sender(ch chan int, count int) {
	for i := 1; i <= count; i++ {
		fmt.Printf("Sender: sending %d\n", i)
		ch <- i // Send data to the channel.
		time.Sleep(time.Millisecond * 300)
	}
	fmt.Println("Sender: closing channel")
	close(ch) // Signal that no more data will be sent.
}

// receiver reads numbers from the channel until it's closed.
func receiver(ch chan int, done chan bool) {
	for num := range ch { // Automatically stops when the channel is closed.
		fmt.Printf("Receiver: received %d\n", num)
		time.Sleep(time.Millisecond * 500)
	}
	fmt.Println("Receiver: no more data, stopping")
	done <- true // Notify main that we're done.
}

func main() {
	ch := make(chan int)    // Unbuffered channel for exchanging data.
	done := make(chan bool) // Used to signal completion.

	// Start sender and receiver goroutines.
	go sender(ch, 5)
	go receiver(ch, done)

	// Wait for the receiver to finish.
	<-done
	fmt.Println("Main: all done!")
}
