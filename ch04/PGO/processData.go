// dramatic_pgo_example.go - Designed to show clear PGO benefits
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

// Interface that will benefit from devirtualization
type Processor interface {
	Process(data []float64) float64
	GetType() string
}

// Hot path processor - used 99% of the time
type FastProcessor struct {
	multiplier float64
}

func (fp FastProcessor) Process(data []float64) float64 {
	sum := 0.0
	// Simple but frequent operation
	for i, v := range data {
		sum += v * fp.multiplier * float64(i+1)
	}
	return sum
}

func (fp FastProcessor) GetType() string {
	return "fast"
}

// Cold path processor - used 1% of the time
type ComplexProcessor struct {
	coefficients []float64
}

func (cp ComplexProcessor) Process(data []float64) float64 {
	if len(cp.coefficients) == 0 {
		cp.coefficients = []float64{1.1, 2.2, 3.3, 4.4, 5.5}
	}

	result := 0.0
	// Expensive operations that are rarely executed
	for i, v := range data {
		for j, coeff := range cp.coefficients {
			result += math.Sin(v*coeff) * math.Cos(float64(i+j))
			result = math.Sqrt(math.Abs(result)) + 0.001
		}
	}
	return result
}

func (cp ComplexProcessor) GetType() string {
	return "complex"
}

// Function with predictable branching (99% one way, 1% the other)
func processData(processor Processor, data []float64, useCache bool) float64 {
	// This branch will be taken 99% of the time
	if useCache {
		// Hot path - simple processing
		return processor.Process(data) * 1.1
	} else {
		// Cold path - expensive processing
		result := processor.Process(data)

		// Additional expensive work (rarely executed)
		for i := 0; i < 100; i++ {
			result = math.Log(math.Abs(result) + 1)
			result *= 1.001
		}

		return result
	}
}

// Function that will benefit from inlining decisions
func calculateMetrics(values []float64) (sum, avg, variance float64) {
	n := float64(len(values))
	if n == 0 {
		return 0, 0, 0
	}

	for _, v := range values {
		sum += v
	}

	avg = sum / n
	for _, v := range values {
		diff := v - avg
		variance += diff * diff
	}
	variance /= n

	return sum, avg, variance
}

// Polymorphic function calls that benefit from devirtualization
func runWorkload(processors []Processor, datasets [][]float64) float64 {
	totalResult := 0.0

	for i := 0; i < 100_000; i++ {
		// 99% of calls use FastProcessor (index 0)
		// 1% of calls use ComplexProcessor (index 1)
		processorIndex := 0
		if rand.Intn(100) == 0 {
			processorIndex = 1
		}

		datasetIndex := i % len(datasets)

		processor := processors[processorIndex]
		useCache := rand.Intn(100) != 0
		result := processData(processor, datasets[datasetIndex], useCache)

		sum, avg, variance := calculateMetrics(datasets[datasetIndex])
		totalResult += result + sum + avg + variance

		// Progress indicator
		if i%10000 == 0 && i > 0 {
			fmt.Printf("Processed %d/100000 items\n", i)
		}
	}

	return totalResult
}

func main() {
	flag.Parse()

	// Enable profiling
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Starting dramatic PGO example...")
	start := time.Now()

	// Create processors (FastProcessor will dominate usage)
	processors := []Processor{
		FastProcessor{multiplier: 1.5},      // Used 99% of the time
		ComplexProcessor{coefficients: nil}, // Used 1% of the time
	}

	// Create test datasets
	datasets := make([][]float64, 5)
	for i := range datasets {
		datasets[i] = make([]float64, 50) // Reasonable size for frequent processing
		for j := range datasets[i] {
			datasets[i][j] = rand.Float64() * 100
		}
	}

	// Run the workload
	result := runWorkload(processors, datasets)

	duration := time.Since(start)

	fmt.Printf("\nCompleted!\n")
	fmt.Printf("Total result: %.2f\n", result)
	fmt.Printf("Execution time: %v\n", duration)
	fmt.Printf("Average time per iteration: %v\n", duration/100000)
}
