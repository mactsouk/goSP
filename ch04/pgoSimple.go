// pgo_example.go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func hotPath() int {
	sum := 0
	for i := 0; i < 1_000_000; i++ {
		sum += i
	}
	return sum
}

func coldPath() int {
	time.Sleep(10 * time.Millisecond)
	return 42
}

func main() {
	rand.Seed(time.Now().UnixNano())
	total := 0

	// Simulate workload: hotPath dominates
	for i := 0; i < 10_000; i++ {
		if rand.Intn(10000) == 0 { // very rare case
			total += coldPath()
		} else {
			total += hotPath()
		}
	}

	fmt.Println("Total:", total)
}
