package main

import (
	"fmt"
	"sync"
	"time"
)

type Request struct {
	Increment int
	Reply     chan int
}

func counterMonitor(requests <-chan Request) {
	count := 0
	for req := range requests {
		count += req.Increment
		req.Reply <- count
	}
	fmt.Println("Monitor exiting")
}

func main() {
	requests := make(chan Request)
	var wg sync.WaitGroup

	numWorkers := 3
	numIncrements := 2

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Go(func() {
			reply := make(chan int)
			for j := 0; j < numIncrements; j++ {
				requests <- Request{Increment: 1, Reply: reply}
				fmt.Printf("Worker %d sees counter = %d\n", i, <-reply)
				time.Sleep(time.Millisecond * 200)
			}
		})
	}

	// Start monitor goroutine separately
	monitorDone := make(chan struct{})
	go func() {
		counterMonitor(requests)
		close(monitorDone) // signal monitor has exited
	}()

	wg.Wait()       // Wait for all workers to finish
	close(requests) // Now it is safe to close requests
	<-monitorDone   // Wait for monitor to exit

	fmt.Println("All workers and monitor finished.")
}
