package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	tasks := []string{"task1", "task2", "task3"}

	for _, t := range tasks {
		wg.Go(func() {
			fmt.Println("Processing:", t)
		})
	}

	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("All tasks completed!")
}
