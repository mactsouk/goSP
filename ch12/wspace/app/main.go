package main

import (
	"fmt"

	// Because of go.work, this imports from ../tslib
	"github.com/mactsouk/tslib"
)

func main() {
	ts1 := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ts2 := []float64{1, 2, 3, 5, 5, 6, 7, 8, 9, 12}

	fmt.Println("--- Time Series Library Test ---")
	fmt.Printf("Series A: %v\n", ts1)
	fmt.Printf("Series B: %v\n\n", ts2)

	// --- Lock-Step Distances ---
	// These compare points i-to-i.
	fmt.Println("1. Lock-Step Measures:")
	euc := tslib.EuclideanDistance(ts1, ts2)
	fmt.Printf("  Euclidean: %.4f\n", euc)
	man := tslib.ManhattanDistance(ts1, ts2)
	fmt.Printf("  Manhattan: %.4f\n", man)
	che := tslib.ChebyshevDistance(ts1, ts2)
	fmt.Printf("  Chebyshev: %.4f\n", che)
	fmt.Println()

	// --- Elastic Distances ---
	// These handle warping and shifting.
	fmt.Println("2. Elastic Measures:")

	// DTW with a window of 2
	dtw := tslib.DTW(ts1, ts2, 2)
	fmt.Printf("  DTW (w=2): %.4f\n", dtw)
	// LCSS with delta (time constraint) = 2 and epsilon (value threshold) = 1.0
	// Returns a distance between 0 (identical) and 1 (no similarity).
	lcss := tslib.LCSS(ts1, ts2, 2, 1.0)
	fmt.Printf("  LCSS (d=2, eps=1.0): %.4f\n", lcss)
	// MPdist with a subsequence window of 3
	mpd := tslib.MPdist(ts1, ts2, 3)
	fmt.Printf("  MPdist (w=3): %.4f\n", mpd)
}
