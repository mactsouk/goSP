// dramatic_pgo_test.go
package main

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkDramaticWorkload(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	// Create processors
	processors := []Processor{
		FastProcessor{multiplier: 1.5},
		ComplexProcessor{coefficients: nil},
	}

	// Create datasets
	datasets := make([][]float64, 5)
	for i := range datasets {
		datasets[i] = make([]float64, 50)
		for j := range datasets[i] {
			datasets[i][j] = rand.Float64() * 100
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// 99% FastProcessor, 1% ComplexProcessor
		processorIndex := 0
		if rand.Intn(100) == 0 {
			processorIndex = 1
		}

		datasetIndex := i % len(datasets)
		processor := processors[processorIndex]

		// 99% cache hit, 1% cache miss
		useCache := rand.Intn(100) != 0

		result := processData(processor, datasets[datasetIndex], useCache)
		calculateMetrics(datasets[datasetIndex])

		// Prevent optimization
		_ = result
	}
}

// Benchmark just the hot interface dispatch
func BenchmarkInterfaceDispatch(b *testing.B) {
	processor := FastProcessor{multiplier: 1.5}
	var p Processor = processor // Force interface dispatch

	data := make([]float64, 50)
	for i := range data {
		data[i] = float64(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := p.Process(data) // Interface call - can be devirtualized
		_ = result
	}
}

// Benchmark the hot function for inlining benefits
func BenchmarkCalculateMetrics(b *testing.B) {
	data := make([]float64, 50)
	for i := range data {
		data[i] = rand.Float64() * 100
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sum, avg, variance := calculateMetrics(data) // Should be inlined
		_ = sum + avg + variance
	}
}
