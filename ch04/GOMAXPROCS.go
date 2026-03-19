package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Simple CPU-intensive work function
func doWork(id int, iterations int) {
	sum := 0
	for i := 0; i < iterations; i++ {
		sum += i * i
	}
	fmt.Printf("Worker %d finished (sum: %d)\n", id, sum)
}

func main() {
	fmt.Printf("Number of CPUs: %d\n", runtime.NumCPU())
	fmt.Printf("Default GOMAXPROCS: %d\n\n", runtime.GOMAXPROCS(0))

	// Test with GOMAXPROCS = 1 (single-threaded)
	fmt.Println("=== Testing with GOMAXPROCS = 1 ===")
	runtime.GOMAXPROCS(1)

	start := time.Now()
	var wg sync.WaitGroup

	// Start 4 goroutines
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			doWork(id, 1000000)
		}(i)
	}

	wg.Wait()
	sThreadT := time.Since(start)
	fmt.Printf("Time with GOMAXPROCS=1: %v\n\n", sThreadT)

	// Test with GOMAXPROCS = NumCPU (multi-threaded)
	maxProcs := runtime.NumCPU()
	fmt.Printf("=== Testing with GOMAXPROCS = %d ===\n", maxProcs)
	runtime.GOMAXPROCS(maxProcs)

	start = time.Now()

	// Start 4 goroutines again
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			doWork(id, 1000000)
		}(i)
	}

	wg.Wait()
	mThreadT := time.Since(start)
	fmt.Printf("Time with GOMAXPROCS=%d: %v\n\n", maxProcs, mThreadT)

	// Show the performance difference
	if sThreadT > mThreadT {
		speedup := float64(sThreadT) / float64(mThreadT)
		fmt.Printf("Multi-threaded was %.2fx faster!\n", speedup)
	} else {
		fmt.Println("Single-threaded was faster!")
	}

	// Show current settings
	fmt.Printf("\nFinal GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
}
