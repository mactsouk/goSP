package main

import (
	"fmt"
	"sync"
)

var counter int

func main() {
	fmt.Println("=== Demonstrating a Race Condition ===")
	counter = 0
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				// Multiple goroutines update `counter` at the same time.
				counter++
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("Final counter value (expected 5000):", counter)
}
