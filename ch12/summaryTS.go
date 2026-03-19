package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

// readGzipTimeSeries reads float64 values from a .gz file.
func readGzipTimeSeries(path string) ([]float64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gr, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	var values []float64
	scanner := bufio.NewScanner(gr)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		val, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float: %s", line)
		}
		values = append(values, val)
	}
	return values, scanner.Err()
}

// 1. PAA (Piecewise Aggregate Approximation)
// Reduces the series to 'targetLen' points by averaging segments.
func PAA(data []float64, targetLen int) []float64 {
	n := len(data)
	if targetLen >= n || targetLen <= 0 {
		return data
	}

	result := make([]float64, targetLen)
	// Segment size (can be a fraction)
	segmentSize := float64(n) / float64(targetLen)

	for i := 0; i < targetLen; i++ {
		// Calculate the start and end indices for this segment in the original array
		start := int(float64(i) * segmentSize)
		end := int(float64(i+1) * segmentSize)

		// Handle rounding edge case at the very end
		if end > n {
			end = n
		}

		// Compute mean of the segment
		sum := 0.0
		count := 0
		for j := start; j < end; j++ {
			sum += data[j]
			count++
		}
		if count > 0 {
			result[i] = sum / float64(count)
		} else {
			// Fallback for very small segments (shouldn't happen if targetLen < n)
			result[i] = data[start]
		}
	}
	return result
}

// 2. Random Sampling
// Selects 'targetLen' random points, preserving chronological order.
func RandomSampling(data []float64, targetLen int) []float64 {
	n := len(data)
	if targetLen >= n || targetLen <= 0 {
		return data
	}

	rand.Seed(time.Now().UnixNano())

	// Use a map to track unique random indices
	indices := make(map[int]bool)
	for len(indices) < targetLen {
		idx := rand.Intn(n)
		indices[idx] = true
	}

	// Extract keys and sort them to preserve time order
	sortedIndices := make([]int, 0, targetLen)
	for idx := range indices {
		sortedIndices = append(sortedIndices, idx)
	}
	sort.Ints(sortedIndices)

	// Build the result
	result := make([]float64, targetLen)
	for i, idx := range sortedIndices {
		result[i] = data[idx]
	}
	return result
}

// 3. Stepping (Stepping/Decimation)
// Selects every k-th element.
func Stepping(data []float64, step int) []float64 {
	if step <= 1 {
		return data
	}
	var result []float64
	for i := 0; i < len(data); i += step {
		result = append(result, data[i])
	}
	return result
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run summarize.go <file.gz> <method> [param]")
		fmt.Println("Methods: paa, random, step")
		fmt.Println("Example: go run summarize.go data.gz paa 100")
		os.Exit(1)
	}

	path := os.Args[1]
	method := os.Args[2]

	// Parse parameter (Target Length for PAA/Random, Step Size for Stepping)
	param := 10 // default
	if len(os.Args) > 3 {
		p, err := strconv.Atoi(os.Args[3])
		if err == nil {
			param = p
		}
	}

	data, err := readGzipTimeSeries(path)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var summary []float64

	switch method {
	case "paa":
		fmt.Printf("Performing PAA to reduce to %d points...\n", param)
		summary = PAA(data, param)
	case "random":
		fmt.Printf("Performing Random Sampling to select %d points...\n", param)
		summary = RandomSampling(data, param)
	case "step":
		fmt.Printf("Performing Stepping with step size %d...\n", param)
		summary = Stepping(data, param)
	default:
		log.Fatalf("Unknown method: %s", method)
	}

	fmt.Printf("Original Length: %d\n", len(data))
	fmt.Printf("Summary Length:  %d\n", len(summary))
	fmt.Println("--- First 10 points of summary ---")
	for i := 0; i < len(summary) && i < 10; i++ {
		fmt.Printf("%.4f ", summary[i])
	}
	fmt.Println()
}
