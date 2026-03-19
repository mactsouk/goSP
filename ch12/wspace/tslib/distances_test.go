package tslib

import (
	"math/rand"
	"testing"
	"time"
)

// helper function to generate random time series
func generateRandomTS(size int) []float64 {
	rand.Seed(time.Now().UnixNano())
	data := make([]float64, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Float64() * 100
	}
	return data
}

// 1. Benchmark Lock-Step Measures
func BenchmarkEuclidean(b *testing.B) {
	ts1 := generateRandomTS(1000)
	ts2 := generateRandomTS(1000)

	b.ResetTimer() // Ignore setup time
	for i := 0; i < b.N; i++ {
		EuclideanDistance(ts1, ts2)
	}
}

func BenchmarkManhattan(b *testing.B) {
	ts1 := generateRandomTS(1000)
	ts2 := generateRandomTS(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ManhattanDistance(ts1, ts2)
	}
}

func BenchmarkChebyshev(b *testing.B) {
	ts1 := generateRandomTS(1000)
	ts2 := generateRandomTS(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ChebyshevDistance(ts1, ts2)
	}
}

// 2. Benchmark Elastic Measures (Expensive!)

func BenchmarkDTW(b *testing.B) {
	// DTW is O(N*M), so 1000 points might be slow.
	// We use a window of 50.
	ts1 := generateRandomTS(1000)
	ts2 := generateRandomTS(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DTW(ts1, ts2, 50)
	}
}

func BenchmarkLCSS(b *testing.B) {
	ts1 := generateRandomTS(1000)
	ts2 := generateRandomTS(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LCSS(ts1, ts2, 50, 0.5)
	}
}

func BenchmarkMPdist(b *testing.B) {
	// MPdist is significantly slower than others.
	// We use a smaller size (200) to keep the benchmark quick enough to run.
	ts1 := generateRandomTS(1000)
	ts2 := generateRandomTS(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Window size of 20
		MPdist(ts1, ts2, 20)
	}
}
