package main

import (
	"fmt"
	"runtime/metrics"
	"time"
)

func main() {
	// Create some activity to populate states
	go func() { time.Sleep(10 * time.Second) }() // Waiting
	go func() {
		for {
			// Busy loop to force "Running" state
			if false {
				break
			}
		}
	}()

	// Define the NEW keys introduced in Go 1.26
	// These specific paths for state breakdown are the new feature.
	descs := []metrics.Sample{
		{Name: "/sched/goroutines/runnable:goroutines"},
		{Name: "/sched/goroutines/running:goroutines"},
		{Name: "/sched/goroutines/waiting:goroutines"},
		{Name: "/sched/goroutines/syscall:goroutines"},
	}

	// Sample the metrics
	metrics.Read(descs)

	// Print the breakdown
	fmt.Println("Go 1.26 Scheduler State Breakdown:")
	for _, sample := range descs {
		// Value is a uint64 count for these new keys
		if sample.Value.Kind() == metrics.KindUint64 {
			fmt.Printf(" - %-35s : %d\n",
				sample.Name, sample.Value.Uint64())
		}
	}
}
