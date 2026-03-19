/*
PROGRAM: Time Series Outlier Detection

DESCRIPTION:
This utility identifies "outliers" in a dataset—time series that
are significantly different from their neighbors.

METHODOLOGY:
The "Outlier Score" for a single file is defined as the Average
Distance to its 'k' Nearest Neighbors. A high score means the file
is isolated.

MODES:
1. Specific Metric (e.g., -metric Euclidean):
   - Computes distances using only that metric.
   - Outlier Score = Raw Average Distance.
   - Files are sorted by this raw score descending.

2. Unified Mode (Default: -metric ALL):
   - Computes Outlier Scores using ALL 6 metrics (Euclidean,
     Manhattan, Chebyshev, DTW, LCSS, MPdist).
   - Because raw distances have incompatible scales (e.g.,
     50,000 vs 0.5), we use Rank Aggregation:
     a. Calculate raw outlier scores for every file for every
        metric.
     b. RANK the files for each metric (1 = Normal, N = Most
        Outlier).
     c. Compute a Weighted Sum of Ranks:
        FinalScore = Σ (Rank_metric * Weight_metric)
        - Lock-Step Weights (Euc, Man, Cheb) = 1.0
        - Elastic Weights (DTW, LCSS, MPdist) = 2.0
   - Files are sorted by this composite rank score.

USAGE:
   go run outliers.go -k 2 -metric ALL [files...]
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/mactsouk/tslib"
)

const (
	WeightLockStep = 1.0
	WeightElastic  = 2.0
)

// Result for a single file in "ALL" mode
type CompositeResult struct {
	Filename   string
	RawScores  map[string]float64 // Map[Metric] -> AvgDist
	Ranks      map[string]int     // Map[Metric] -> Rank
	FinalScore float64
}

// Result for single metric mode
type SimpleResult struct {
	Filename string
	Score    float64
}

func main() {
	// --- Parameters ---
	k := flag.Int("k", 1,
		"Number of neighbors to average distance over")
	metric := flag.String("metric", "ALL",
		"Distance metric (Euclidean, DTW, etc.) or 'ALL'")
	window := flag.Int("w", 20,
		"Window size for elastic metrics")
	targetLen := flag.Int("L", 0,
		"Target length (0 = disabled)")
	summaryMethod := flag.String("S", "paa",
		"Summarization method")
	flag.Parse()

	files := flag.Args()
	if len(files) < *k+1 {
		log.Fatal("Error: Not enough files to compute " +
			"neighbors for outlier detection.")
	}

	// --- 0. Print Active Parameters ---
	fmt.Println("--------------------------------------------")
	fmt.Println("       Time Series Outlier Detection")
	fmt.Println("--------------------------------------------")
	fmt.Printf("Metric: %s\n", *metric)
	if *metric == "ALL" || isElastic(*metric) {
		fmt.Printf("Config: Window=%d (Elastic only)\n",
			*window)
	}
	fmt.Printf("Config: Neighbors(k)=%d\n", *k)
	fmt.Printf("Summary: Length=%d (0=Full), Method=%s\n",
		*targetLen, *summaryMethod)
	fmt.Printf("Input:  %d files\n", len(files))
	fmt.Println("--------------------------------------------")

	// 1. Load Data
	fmt.Print("Loading data... ")
	data := make(map[string][]float64)
	for _, f := range files {
		data[f] = loadAndProcess(f, *targetLen,
			*summaryMethod)
	}
	fmt.Println("Done.")

	// 2. Validate Data Consistency
	if len(data) == 0 {
		log.Fatal("Error: No data loaded.")
	}

	var expectedLen int
	firstFile := true
	for fname, ts := range data {
		if len(ts) == 0 {
			log.Fatalf("Error: Empty time series (%s)",
				fname)
		}
		if firstFile {
			expectedLen = len(ts)
			firstFile = false
		} else if len(ts) != expectedLen {
			log.Fatalf("Error: Inconsistent time series "+
				"length (%s): expected %d, got %d",
				fname, expectedLen, len(ts))
		}
	}
	fmt.Printf("Validated: All %d time series have "+
		"length %d\n", len(data), expectedLen)

	// Dispatch based on mode
	if *metric == "ALL" {
		runCompositeMode(files, data, *k, *window)
	} else {
		runSingleMode(files, data, *metric, *k, *window)
	}
}

// --- Mode 1: Compute ALL metrics and Aggregate Ranks ---
func runCompositeMode(files []string,
	data map[string][]float64, k, window int) {
	metrics := []string{"Euclidean", "Manhattan", "Chebyshev",
		"DTW", "LCSS", "MPdist"}

	// Store raw scores for everyone
	// results[filename].RawScores["Euclidean"] = 50.4
	results := make(map[string]*CompositeResult)
	for _, f := range files {
		results[f] = &CompositeResult{
			Filename:  f,
			RawScores: make(map[string]float64),
			Ranks:     make(map[string]int),
		}
	}

	fmt.Println("Calculating ALL metrics (Parallel)...")

	// Calculate one metric at a time to simplify ranking
	for _, m := range metrics {
		scores := computeScoresForMetric(files, data, m, k,
			window)

		// Store raw scores
		for _, s := range scores {
			results[s.Filename].RawScores[m] = s.Score
		}

		// Rank: Ascending Score = Least Outlier = Rank 1
		// High Score = Most Outlier = Rank N
		sort.Slice(scores, func(i, j int) bool {
			return scores[i].Score < scores[j].Score
		})

		for rank, s := range scores {
			// Rank 1 = smallest distance (most normal)
			// Rank N = largest distance (most outlier)
			results[s.Filename].Ranks[m] = rank + 1
		}
	}

	// Compute Final Weighted Score
	var finalBox []*CompositeResult
	for _, r := range results {
		wSum := 0.0

		wSum += float64(r.Ranks["Euclidean"]) *
			WeightLockStep
		wSum += float64(r.Ranks["Manhattan"]) *
			WeightLockStep
		wSum += float64(r.Ranks["Chebyshev"]) *
			WeightLockStep

		wSum += float64(r.Ranks["DTW"]) * WeightElastic
		wSum += float64(r.Ranks["LCSS"]) * WeightElastic
		wSum += float64(r.Ranks["MPdist"]) * WeightElastic

		r.FinalScore = wSum
		finalBox = append(finalBox, r)
	}

	// Sort by Final Score Descending (Biggest Outlier first)
	sort.Slice(finalBox, func(i, j int) bool {
		return finalBox[i].FinalScore >
			finalBox[j].FinalScore
	})

	// Print
	fmt.Println("\n--- Composite Outlier Report " +
		"(Weighted Ranks) ---")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Rank\tFilename\tComposite Score\t"+
		"Top Metric Contributors (Rank)")
	fmt.Fprintln(w, "----\t--------\t---------------\t"+
		"------------------------------")

	for i, r := range finalBox {
		// Create a quick summary of where this file
		// ranked highest e.g. "DTW(#8) MPdist(#8)"
		contribs := getTopContributors(r.Ranks)
		fmt.Fprintf(w, "%d\t%s\t%.1f\t%s\n", i+1,
			r.Filename, r.FinalScore, contribs)
	}
	w.Flush()
}

// --- Mode 2: Standard Single Metric ---
func runSingleMode(files []string, data map[string][]float64,
	metric string, k, window int) {
	scores := computeScoresForMetric(files, data, metric, k,
		window)

	// Sort Descending (High dist = Outlier)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	fmt.Printf("\n--- Outlier Detection Report "+
		"(Metric: %s) ---\n", metric)
	fmt.Println("Rank  Filename        " +
		"Outlier Score (Avg Dist)")
	fmt.Println("----  --------        " +
		"------------------------")
	for i, s := range scores {
		fmt.Printf("%-4d  %-14s  %.4f\n", i+1, s.Filename,
			s.Score)
	}
}

// --- Core Logic: Compute Avg Distance to k-NN ---
func computeScoresForMetric(files []string,
	data map[string][]float64, metric string,
	k, window int) []SimpleResult {

	// Pre-allocate results slice
	results := make([]SimpleResult, len(files))
	var wg sync.WaitGroup

	for i, currentFile := range files {
		wg.Add(1)
		go func(idx int, fname string) {
			defer wg.Done()

			currentTS := data[fname]

			// Pre-allocate distances slice
			distances := make([]float64, 0,
				len(files)-1)

			for _, otherFile := range files {
				if fname == otherFile {
					continue
				}

				otherTS := data[otherFile]
				d := 0.0
				switch strings.ToLower(metric) {
				case "dtw":
					d = tslib.DTW(currentTS, otherTS,
						window)
				case "mpdist":
					d = tslib.MPdist(currentTS,
						otherTS, window)
				case "lcss":
					d = tslib.LCSS(currentTS, otherTS,
						window, 0.5)
				case "manhattan":
					d = tslib.ManhattanDistance(
						currentTS, otherTS)
				case "chebyshev":
					d = tslib.ChebyshevDistance(
						currentTS, otherTS)
				default:
					d = tslib.EuclideanDistance(
						currentTS, otherTS)
				}
				distances = append(distances, d)
			}

			sort.Float64s(distances)

			sum := 0.0
			for j := 0; j < k && j < len(distances); j++ {
				sum += distances[j]
			}
			avg := sum / float64(k)

			// Write directly to pre-allocated index
			results[idx] = SimpleResult{fname, avg}
		}(i, currentFile)
	}
	wg.Wait()
	return results
}

// --- Helpers ---

func getTopContributors(ranks map[string]int) string {
	// Helper to print which metrics ranked this file highest
	// Sort metrics by rank descending
	type kv struct {
		Key string
		Val int
	}
	var sorted []kv
	for k, v := range ranks {
		sorted = append(sorted, kv{k, v})
	}
	// Descending by rank
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Val > sorted[j].Val
	})

	// Return top 3
	var parts []string
	for i := 0; i < 3 && i < len(sorted); i++ {
		parts = append(parts,
			fmt.Sprintf("%s(#%d)", sorted[i].Key,
				sorted[i].Val))
	}
	return strings.Join(parts, " ")
}

func isElastic(m string) bool {
	m = strings.ToLower(m)
	return m == "dtw" || m == "mpdist" || m == "lcss"
}

func loadAndProcess(path string, length int,
	method string) []float64 {
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
