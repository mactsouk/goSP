package main

import (
	"fmt"
	"runtime"
	"time"
)

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MB, TotalAlloc = %v MB, NumGoroutine = %v\n",
		m.Alloc/1024/1024, m.TotalAlloc/1024/1024, runtime.NumGoroutine())
}

func main() {
	for i := 0; i < 50; i++ { // Launch 50 goroutines quickly
		go func(id int) {
			data := make([]byte, 5*1024*1024) // 5 MB slice per goroutine
			for {
				_ = data // Keep reference alive
				time.Sleep(500 * time.Millisecond)
			}
		}(i)
		time.Sleep(50 * time.Millisecond) // Slight delay to stagger allocations
		if i%5 == 0 {
			printMemUsage()
		}
	}

	fmt.Println("All goroutines launched–memory continues to grow.")
	time.Sleep(5 * time.Second)
}
