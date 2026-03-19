package main

import (
	"fmt"
	"sync"
)

var counter int
var mu sync.Mutex

func main() {
	fmt.Println("=== Fixed Race Condition ===")
	counter = 0
	var wg sync.WaitGroup

	// Launch 5 goroutines that increment the same counter safely.
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				// Protect access to counter with a mutex.
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Final counter value (expected 5000):", counter)
}
