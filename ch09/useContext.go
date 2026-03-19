package main

import (
	"context"
	"fmt"
	"time"
)

func doWork(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Work canceled:", ctx.Err())
			return
		default:
			fmt.Print("Working... ")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main() {
	// Create a context that automatically cancels after 2 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Always call cancel to release resources

	fmt.Println("Starting work")
	doWork(ctx)
	fmt.Println("Finished")
}
