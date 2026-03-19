package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/mactsouk/tslib" // Change to your actual username
)

// Command Line Flags
var (
	targetLen     int
	summaryMethod string
	window        int
	epsilon       float64
)

func main() {
	flag.IntVar(&targetLen, "L", 0, "Target length for summarization (0 = disabled)")
	flag.StringVar(&summaryMethod, "S", "paa", "Summarization method: 'paa', 'random', 'step'")
	flag.IntVar(&window, "w", 20, "Base window size for elastic measures")
	flag.Float64Var(&epsilon, "e", 0.0, "Epsilon threshold for LCSS")
	flag.Parse()

	files := flag.Args()
	if len(files) < 2 {
		fmt.Println("Usage: go run nn_finder.go -L 50 [files...]")
		os.Exit(1)
	}

	// 1. Read and Process Data
	dataMap := make(map[string][]float64)
	var filenames []string

	for _, f := range files {
		raw, err := tslib.ReadGzipTimeSeries(f)
		if err != nil {
			log.Printf("Skipping %s: %v", f, err)
			continue
		}

		// Apply Summarization
		var processed []float64
		if targetLen > 0 {
			switch strings.ToLower(summaryMethod) {
			case "paa":
				processed = tslib.PAA(raw, targetLen)
			case "random":
				processed = tslib.RandomSampling(raw, targetLen)
			case "step":
				processed = tslib.Stepping(raw, targetLen)
			default:
				log.Fatalf("Unknown method: %s", summaryMethod)
			}
		} else {
			processed = raw
		}
		dataMap[f] = processed
		filenames = append(filenames, f)
	}

	if len(filenames) < 2 {
		log.Fatal("Not enough valid files.")
	}

	// 2. Find Nearest Neighbors
	metrics := []string{"Euclidean", "Manhattan", "Chebyshev", "DTW", "LCSS", "MPdist"}

	for _, queryName := range filenames {
		query := dataMap[queryName]

		bestDist := make(map[string]float64)
		bestMatch := make(map[string]string)

		for _, m := range metrics {
			bestDist[m] = math.Inf(1)
			bestMatch[m] = "-"
		}

		for _, candidateName := range filenames {
			if queryName == candidateName {
				continue
			}
			candidate := dataMap[candidateName]

			// Auto-correct Window and Epsilon
			lenDiff := abs(len(query) - len(candidate))
			effectiveWindow := window
			if lenDiff > effectiveWindow {
				effectiveWindow = lenDiff + window
			}

			effectiveEpsilon := epsilon
			if effectiveEpsilon == 0.0 {
				effectiveEpsilon = getStdDev(candidate) * 0.2
			}

			dists := make(map[string]float64)
			dists["Euclidean"] = tslib.EuclideanDistance(query, candidate)
			dists["Manhattan"] = tslib.ManhattanDistance(query, candidate)
			dists["Chebyshev"] = tslib.ChebyshevDistance(query, candidate)
			dists["DTW"] = tslib.DTW(query, candidate, effectiveWindow)
			dists["LCSS"] = tslib.LCSS(query, candidate, effectiveWindow, effectiveEpsilon)
			dists["MPdist"] = tslib.MPdist(query, candidate, effectiveWindow)

			for _, m := range metrics {
				if dists[m] < bestDist[m] {
					bestDist[m] = dists[m]
					bestMatch[m] = candidateName
				}
			}
		}

		printLineReport(queryName, bestMatch, bestDist, metrics)
	}
}

// Updated printing function for single-line output
func printLineReport(query string, matches map[string]string, scores map[string]float64, order []string) {
	shortQuery := filepath.Base(query)
	fmt.Printf("[%s] ", shortQuery)

	for i, m := range order {
		shortNeighbor := filepath.Base(matches[m])
		// Format: Metric: Neighbor (Dist)
		fmt.Printf("%s: %s (%.2f)", m, shortNeighbor, scores[m])

		// Add separator unless it's the last item
		if i < len(order)-1 {
			fmt.Print(" | ")
		}
	}
	fmt.Println() // Newline at the end
}

// --- Helpers ---

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func getStdDev(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	mean := sum / float64(len(data))
	sqDiff := 0.0
	for _, v := range data {
		diff := v - mean
		sqDiff += diff * diff
	}
	return math.Sqrt(sqDiff / float64(len(data)))
}
