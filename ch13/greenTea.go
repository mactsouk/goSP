package main

import (
	"fmt"
	"runtime"
	"time"
)

// Node represents a single object in our graph.
// We add padding (Payload) to ensure each Node takes up significant space,
// forcing the CPU to fetch new cache lines when traversing the list.
type Node struct {
	Next    *Node
	Payload [64]byte
}

// Configuration
const (
	// LiveSetSize: We create 5 million persistent objects.
	// The GC MUST scan all of these every single time it runs.
	LiveSetSize = 5_000_000

	// ChurnAmount: We allocate 50 million temporary objects.
	// This fills the heap rapidly, triggering the GC over and over.
	ChurnAmount = 50_000_000
)

func main() {
	fmt.Println("1. Setup: Building the massive Live Set (Linked List)...")

	// Construct a giant linked list.
	// This structure is "pointer-heavy" and hard to scan.
	root := &Node{}
	current := root
	for i := 0; i < LiveSetSize; i++ {
		newNode := &Node{}
		current.Next = newNode
		current = newNode
	}

	// Force a GC now to clean up setup artifacts and stabilize the heap.
	runtime.GC()

	fmt.Println("2. Benchmark: Starting allocation churn...")
	start := time.Now()

	// The "Churn" Loop
	// We rapidly allocate and discard objects. This creates "memory pressure."
	// Every time the heap fills up, the GC wakes up.
	// To free this memory, the GC must first prove that the 'root' list
	// is still alive by traversing all 5,000,000 pointers.
	for i := 0; i < ChurnAmount; i++ {
		_ = &Node{}
	}

	duration := time.Since(start)

	// Collect final statistics
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	fmt.Println("\n--- Results ---")
	fmt.Printf("Total Wall Time:    %v\n", duration)
	fmt.Printf("Total GC Cycles:    %d\n", stats.NumGC)
	fmt.Printf("Total GC Pause:     %v\n", time.Duration(stats.PauseTotalNs))
	fmt.Printf("Avg Time per Cycle: %v\n", duration/time.Duration(stats.NumGC))

	// Keep 'root' alive so the compiler doesn't optimize the list away
	_ = root
}
