package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	tasks := []string{"task1", "task2", "task3"}

	for _, t := range tasks {
		// Manually increment the counter
		wg.Add(1)
		go func(task string) {
			// Remember to call Done()
			defer wg.Done()
			fmt.Println("Processing:", task)
		}(t)
	}

	wg.Wait()
	fmt.Println("All tasks completed!")
}
