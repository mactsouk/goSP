package main

import (
	"fmt"
	"runtime"
	"time"
)

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MB, TotalAlloc = %v MB, NumGC = %v\n",
		m.Alloc/1024/1024, m.TotalAlloc/1024/1024, m.NumGC)
}

type Data struct {
	payload [1024 * 1024]byte // 1 MB per instance
}

// Common memory leak pattern: goroutine + channel
func leakyFunction(data *Data) {
	ch := make(chan bool) // Unbuffered channel that never receives
	go func() {
		ch <- true // This goroutine blocks forever, keeping 'data' alive
	}()
	// Function returns but goroutine is stuck, keeping reference to data
}


func main() {
	fmt.Println("Memory leak example (stuck goroutines)...")
	
	printMemUsage() // Baseline
	
	for i := 0; i < 10; i++ {
		d := &Data{}
		leakyFunction(d)
		printMemUsage()
		time.Sleep(500 * time.Millisecond)
	}
	
	fmt.Printf("# of goroutines: %d (should be 11: main + 10 stuck)\n",
		runtime.NumGoroutine())
	
	// Force multiple GC cycles to demonstrate leak
	fmt.Println("\nForcing 3 GC cycles to try to clean up...")
	for i := 0; i < 3; i++ {
		runtime.GC()
		time.Sleep(100 * time.Millisecond)
		printMemUsage()
	}
}
