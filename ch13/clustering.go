package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/mactsouk/tslib"
)

func main() {
	// --- Command Line Flags ---
	k := flag.Int("k", 2, "Number of clusters")
	iter := flag.Int("iter", 10, "Maximum iterations")
	targetLen := flag.Int("L", 0,
		"Target length (Required for correct centroid averaging)")
	summaryMethod := flag.String("S", "paa",
		"Summarization method")
	seed := flag.Int64("seed", 0,
		"Random seed (0 for current time)")
	flag.Parse()

	files := flag.Args()
	if len(files) < *k {
		log.Fatal("Error: Number of files must be greater " +
			"than number of clusters (k).")
	}

	// --- 0. Print Active Parameters ---
	fmt.Println("--------------------------------------------")
	fmt.Println("       K-Means Time Series Clustering")
	fmt.Println("--------------------------------------------")
	fmt.Printf("Config: Clusters(k)=%d, MaxIter=%d\n",
		*k, *iter)
	fmt.Printf("Summary: Length=%d (0=Full), Method=%s\n",
		*targetLen, *summaryMethod)
	fmt.Printf("Input:  %d files\n", len(files))
	fmt.Println("--------------------------------------------")

	// 1. Load All Data (Serial I/O)
	fmt.Print("Loading data... ")
	var data [][]float64
	var fileNames []string
	for _, f := range files {
		ts := loadAndProcess(f, *targetLen, *summaryMethod)
		data = append(data, ts)
		fileNames = append(fileNames, f)
	}
	fmt.Println("Done.")

	// 2. Validate Data Consistency
	if len(data) == 0 {
		log.Fatal("Error: No data loaded.")
	}
	expectedLen := len(data[0])
	for i, ts := range data {
		if len(ts) != expectedLen {
			log.Fatalf("Error: Inconsistent time series "+
				"length at index %d (%s): "+
				"expected %d, got %d",
				i, fileNames[i], expectedLen, len(ts))
		}
		if len(ts) == 0 {
			log.Fatalf("Error: Empty time series "+
				"at index %d (%s)",
				i, fileNames[i])
		}
	}
	fmt.Printf("Validated: All %d time series have length %d\n",
		len(data), expectedLen)

	// 3. Initialize Random Number Generator
	var rng *rand.Rand
	if *seed == 0 {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		rng = rand.New(rand.NewSource(*seed))
		fmt.Printf("Using random seed: %d\n", *seed)
	}

	// 4. Initialize Centroids (Random initialization)
	centroids := make([][]float64, *k)
	perm := rng.Perm(len(data))
	for i := 0; i < *k; i++ {
		centroids[i] = make([]float64, expectedLen)
		copy(centroids[i], data[perm[i]])
	}

	assignments := make([]int, len(data))
	for i := range assignments {
		assignments[i] = -1 // Initialize to unassigned
	}

	// 5. K-Means Loop
	fmt.Println("--------------------------------------------")
	for it := 0; it < *iter; it++ {
		changes := 0

		// --- Step A: Assignment (Serial - no race conditions)
		for i := 0; i < len(data); i++ {
			ts := data[i]
			bestCluster := -1
			minDist := math.Inf(1)

			for cIdx, center := range centroids {
				d := tslib.EuclideanDistance(ts, center)
				if d < minDist {
					minDist = d
					bestCluster = cIdx
				}
			}

			if assignments[i] != bestCluster {
				assignments[i] = bestCluster
				changes++
			}
		}

		// --- Step B: Update Centroids ---
		emptyClusterCount := 0
		for cIdx := 0; cIdx < *k; cIdx++ {
			sum := make([]float64, expectedLen)
			count := 0

			for i, a := range assignments {
				if a == cIdx {
					for j := 0; j < expectedLen; j++ {
						sum[j] += data[i][j]
					}
					count++
				}
			}

			if count > 0 {
				// Update centroid with average
				for j := range sum {
					centroids[cIdx][j] = sum[j] /
						float64(count)
				}
			} else {
				// Handle empty cluster: reinitialize
				// with random data point
				emptyClusterCount++
				randomIdx := rng.Intn(len(data))
				copy(centroids[cIdx], data[randomIdx])
				fmt.Printf("  Warning: Cluster %d empty, "+
					"reinitialized with random "+
					"point\n", cIdx)
			}
		}

		fmt.Printf("Iteration %d: %d reassignment(s)",
			it+1, changes)
		if emptyClusterCount > 0 {
			fmt.Printf(" (%d empty cluster(s) "+
				"reinitialized)", emptyClusterCount)
		}
		fmt.Println()

		if changes == 0 {
			fmt.Println("Converged (no changes).")
			break
		}
	}

	// 6. Print Results
	fmt.Println("--------------------------------------------")
	fmt.Printf("--- Final Clusters (k=%d) ---\n", *k)
	for cIdx := 0; cIdx < *k; cIdx++ {
		var members []string
		for i, a := range assignments {
			if a == cIdx {
				// Use filepath.Base to get filename only
				members = append(members,
					filepath.Base(fileNames[i]))
			}
		}
		// Print all members on one line
		if len(members) > 0 {
			fmt.Printf("Cluster %d (%d members): %s\n",
				cIdx, len(members),
				strings.Join(members, ", "))
		} else {
			fmt.Printf("Cluster %d (0 members): <empty>\n",
				cIdx)
		}
	}
	fmt.Println("--------------------------------------------")
}

// loadAndProcess helper
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
