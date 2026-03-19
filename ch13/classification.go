/*
PROGRAM: Leave-One-Out Cross-Validation (LOOCV) for Time Series Similarity

DESCRIPTION:
This utility performs a robust similarity search across a dataset of time series files.
It identifies the single "Nearest Neighbor" for every input file by comparing it against
all other files in the set.

METHODOLOGY:

1.  Z-Normalization (Preprocessing):
    Before any distance calculation, every time series is Z-normalized:
        x' = (x - μ) / σ
    This transforms the data to have a mean of 0 and a standard deviation of 1.
    This step is critical because it removes "amplitude scaling" and "offset" differences,
    allowing us to compare the *shape* of the signals rather than their absolute raw values.

2.  Leave-One-Out Strategy:
    The program iterates through the file list one by one. For each iteration, the
    current file is treated as the "Query" (unknown), and it is compared against the
    remaining n-1 files (the "Candidates").

3.  Weighted Rank Aggregation (The Metric):
    Instead of relying on a single distance metric (like Euclidean), we calculate
    6 different distances for every pair:
      - Lock-Step: Euclidean, Manhattan, Chebyshev
      - Elastic:   DTW, LCSS, MPdist

    Because these metrics have vastly different scales (e.g., Manhattan might be 100,000
    while MPdist is 0.5), we cannot sum them directly. Instead, we RANK the candidates.
    If File B is the closest to File A using Euclidean distance, it gets Rank 1 for Euclidean.

    The final "Similarity Score" is a weighted sum of these ranks:
       Score = Σ (Rank_i * Weight_i)

    Weights are assigned based on importance:
       - Lock-Step Metrics: Weight = 1.0
       - Elastic Metrics:   Weight = 2.0 (Higher priority to capture shape warping)

    The candidate with the LOWEST score is declared the Unified Nearest Neighbor.

USAGE:
    go run loocv.go [flags] [file1.gz] [file2.gz] ...
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/mactsouk/tslib"
)

// Weights for the metrics
const (
	WeightLockStep = 1.0
	WeightElastic  = 2.0
)

type DistanceResult struct {
	Filename   string
	Values     map[string]float64
	Ranks      map[string]int
	FinalScore float64
}

// Struct to store the final decision for each file
type MatchDecision struct {
	Query    string
	Neighbor string
	Score    float64
	Ranks    map[string]int // Store ranks for display
}

func main() {
	// --- Parameters ---
	window := flag.Int("w", 20, "Window size for elastic metrics")
	targetLen := flag.Int("L", 0, "Target length for summarization (0=disabled)")
	summaryMethod := flag.String("S", "paa", "Summarization method")
	flag.Parse()

	files := flag.Args()
	if len(files) < 2 {
		log.Fatal("Need at least 2 files to perform comparisons.")
	}

	// --- 0. Print Active Parameters ---
	fmt.Println("------------------------------------------------")
	fmt.Println("       LOOCV Time Series Similarity")
	fmt.Println("------------------------------------------------")
	fmt.Printf("Config: Window=%d, SummaryLen=%d, Method=%s\n", *window, *targetLen, *summaryMethod)
	fmt.Printf("Input:  %d files\n", len(files))
	fmt.Println("------------------------------------------------")

	// 1. Load and Z-Normalize Data
	fmt.Print("Loading and processing data... ")
	dataMap := make(map[string][]float64)
	for _, f := range files {
		raw := loadAndProcess(f, *targetLen, *summaryMethod)
		norm := zNormalize(raw)
		dataMap[f] = norm
	}
	fmt.Println("Done.")

	metrics := []string{"Euclidean", "Manhattan", "Chebyshev", "DTW", "LCSS", "MPdist"}

	// Store all findings to analyze mutual groups later
	var decisions []MatchDecision
	// Map to quickly look up who 'A' picked (A -> B)
	lookup := make(map[string]string)

	fmt.Println("Calculating distances (Parallel)...")

	// 2. Leave-One-Out Loop
	for _, queryFile := range files {
		queryTS := dataMap[queryFile]

		// Parallel comparison for this query
		var results []*DistanceResult
		var mu sync.Mutex
		var wg sync.WaitGroup
		sem := make(chan struct{}, runtime.NumCPU())

		for _, candidateFile := range files {
			if queryFile == candidateFile {
				continue
			}

			wg.Add(1)
			go func(candFile string) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				candTS := dataMap[candFile]
				dists := make(map[string]float64)

				dists["Euclidean"] = tslib.EuclideanDistance(queryTS, candTS)
				dists["Manhattan"] = tslib.ManhattanDistance(queryTS, candTS)
				dists["Chebyshev"] = tslib.ChebyshevDistance(queryTS, candTS)
				dists["DTW"] = tslib.DTW(queryTS, candTS, *window)
				dists["LCSS"] = tslib.LCSS(queryTS, candTS, *window, 0.5)
				dists["MPdist"] = tslib.MPdist(queryTS, candTS, *window)

				res := &DistanceResult{
					Filename: candFile,
					Values:   dists,
					Ranks:    make(map[string]int),
				}

				mu.Lock()
				results = append(results, res)
				mu.Unlock()
			}(candidateFile)
		}
		wg.Wait()

		// 3. Rank and Score
		for _, m := range metrics {
			sort.Slice(results, func(a, b int) bool {
				return results[a].Values[m] < results[b].Values[m]
			})
			for rank, res := range results {
				res.Ranks[m] = rank + 1
			}
		}

		for _, res := range results {
			score := 0.0
			score += float64(res.Ranks["Euclidean"]) * WeightLockStep
			score += float64(res.Ranks["Manhattan"]) * WeightLockStep
			score += float64(res.Ranks["Chebyshev"]) * WeightLockStep
			score += float64(res.Ranks["DTW"]) * WeightElastic
			score += float64(res.Ranks["LCSS"]) * WeightElastic
			score += float64(res.Ranks["MPdist"]) * WeightElastic
			res.FinalScore = score
		}

		// Find Winner
		sort.Slice(results, func(a, b int) bool {
			return results[a].FinalScore < results[b].FinalScore
		})
		best := results[0]

		decisions = append(decisions, MatchDecision{
			Query:    queryFile,
			Neighbor: best.Filename,
			Score:    best.FinalScore,
			Ranks:    best.Ranks,
		})
		lookup[queryFile] = best.Filename
	}

	// 4. Detailed Report
	fmt.Println("\n--- Nearest Neighbor Report ---")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	// Header
	fmt.Fprintln(w, "Query\tNeighbor\tScore\tRank Breakdown")
	fmt.Fprintln(w, "-----\t--------\t-----\t--------------")

	for _, d := range decisions {
		qShort := filepath.Base(d.Query)
		nShort := filepath.Base(d.Neighbor)

		// Build the rank string
		rankStr := fmt.Sprintf("LockStep: Euc(#%d) Man(#%d) Cheb(#%d) | Elastic: DTW(#%d) LCSS(#%d) MPdist(#%d)",
			d.Ranks["Euclidean"], d.Ranks["Manhattan"], d.Ranks["Chebyshev"],
			d.Ranks["DTW"], d.Ranks["LCSS"], d.Ranks["MPdist"])

		fmt.Fprintf(w, "%s\t%s\t%.1f\t%s\n", qShort, nShort, d.Score, rankStr)
	}
	w.Flush()

	// 5. Identify Mutual Groups (Symmetric NN)
	fmt.Println("\n--- Mutual Nearest Neighbors (Symmetric Groups) ---")
	printed := make(map[string]bool)

	foundGroup := false
	for _, d := range decisions {
		fileA := d.Query
		fileB := d.Neighbor

		if lookup[fileB] == fileA {
			key1 := fileA + fileB
			key2 := fileB + fileA
			if !printed[key1] && !printed[key2] {
				fmt.Printf("%s  <-->  %s\n", filepath.Base(fileA), filepath.Base(fileB))
				printed[key1] = true
				printed[key2] = true
				foundGroup = true
			}
		}
	}
	if !foundGroup {
		fmt.Println("No symmetric pairs found.")
	}
	fmt.Println()
}

// --- Helpers ---

// zNormalize calculates (x - mean) / stdDev
func zNormalize(data []float64) []float64 {
	if len(data) == 0 {
		return data
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	mean := sum / float64(len(data))
	sqDiffSum := 0.0
	for _, v := range data {
		diff := v - mean
		sqDiffSum += diff * diff
	}
	stdDev := math.Sqrt(sqDiffSum / float64(len(data)))
	if stdDev == 0 {
		stdDev = 1.0
	}
	result := make([]float64, len(data))
	for i, v := range data {
		result[i] = (v - mean) / stdDev
	}
	return result
}

func loadAndProcess(path string, length int, method string) []float64 {
	raw, err := tslib.ReadGzipTimeSeries(path)
	if err != nil {
		log.Fatalf("Error reading %s: %v", path, err)
	}
	if length <= 0 {
		return raw
	}
	switch strings.ToLower(method) {
	case "paa":
		return tslib.PAA(raw, length)
	case "random":
		return tslib.RandomSampling(raw, length)
	case "step":
		return tslib.Stepping(raw, length)
	default:
		return raw
	}
}
