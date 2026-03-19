package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// --- Command Line Flags ---
var (
	windowSize    int
	doZNorm       bool
	epsilon       float64
	targetLength  int    // -L
	summaryMethod string // -S
)

func main() {
	flag.IntVar(&windowSize, "w", 50, "Sliding window size (applied AFTER summarization)")
	flag.BoolVar(&doZNorm, "z", false, "Apply z-normalization")
	flag.Float64Var(&epsilon, "e", 0.0, "Epsilon matching threshold for LCSS")

	// New Flags
	flag.IntVar(&targetLength, "L", 0, "Target length for summarization (0 = disabled)")
	flag.StringVar(&summaryMethod, "S", "paa", "Summarization method: 'paa', 'random', 'step'")

	flag.Parse()

	files := flag.Args()
	if len(files) < 2 {
		fmt.Println("Usage: go run main.go -L 100 -S paa -w 10 [files...]")
		os.Exit(1)
	}

	// Read and potentially summarize all series
	var series [][]float64
	for _, f := range files {
		data, err := readTimeSeries(f)
		if err != nil {
			log.Fatalf("Error reading %s: %v", f, err)
		}

		// Apply Summarization if -L is set
		if targetLength > 0 && targetLength < len(data) {
			switch strings.ToLower(summaryMethod) {
			case "paa":
				data = PAA(data, targetLength)
			case "random":
				data = RandomSampling(data, targetLength)
			case "step":
				// For stepping, we calculate step size to approximate target length
				stepSize := len(data) / targetLength
				if stepSize < 1 {
					stepSize = 1
				}
				data = Stepping(data, stepSize)
			default:
				log.Fatalf("Unknown summary method: %s", summaryMethod)
			}
		}
		series = append(series, data)
	}

	query := series[0]
	fmt.Printf("Query: %s (Processed Len: %d)\n", files[0], len(query))
	if targetLength > 0 {
		fmt.Printf("Summarization: %s -> Length %d\n", strings.ToUpper(summaryMethod), targetLength)
	}
	fmt.Printf("Config: Window=%d, Z-Norm=%v\n\n", windowSize, doZNorm)

	header := fmt.Sprintf("%-20s %-12s %-12s %-12s %-12s %-12s %-12s",
		"Target", "Euclidean", "Manhattan", "Chebyshev", "DTW", "LCSS", "MPdist")
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", len(header)+5))

	for i := 1; i < len(series); i++ {
		target := series[i]

		// 1. Lock-step measures
		minLen := min(len(query), len(target))
		qTrim, tTrim := query[:minLen], target[:minLen]

		if doZNorm {
			qTrim = zNormalize(qTrim)
			tTrim = zNormalize(tTrim)
		}

		euc := EuclideanDistance(qTrim, tTrim)
		man := ManhattanDistance(qTrim, tTrim)
		che := ChebyshevDistance(qTrim, tTrim)

		// 2. Elastic measures
		qFull, tFull := query, target
		if doZNorm {
			qFull = zNormalize(query)
			tFull = zNormalize(target)
		}

		dtw := DTW(qFull, tFull, windowSize*2)

		eps := epsilon
		if eps == 0.0 {
			stats := computeStats(tFull)
			eps = 0.2 * stats.StdDev
		}
		lcss := LCSS(qFull, tFull, windowSize, eps)

		// 3. MPdist
		mpd := MPdist(query, target, windowSize, doZNorm)

		fmt.Printf("%-20s %-12.4f %-12.4f %-12.4f %-12.4f %-12.4f %-12.4f\n",
			files[i], euc, man, che, dtw, lcss, mpd)
	}
}

// --- Summarization Functions ---

func PAA(data []float64, targetLen int) []float64 {
	n := len(data)
	if targetLen >= n || targetLen <= 0 {
		return data
	}

	result := make([]float64, targetLen)
	segmentSize := float64(n) / float64(targetLen)

	for i := 0; i < targetLen; i++ {
		start := int(float64(i) * segmentSize)
		end := int(float64(i+1) * segmentSize)
		if end > n {
			end = n
		}

		sum := 0.0
		count := 0
		for j := start; j < end; j++ {
			sum += data[j]
			count++
		}
		if count > 0 {
			result[i] = sum / float64(count)
		} else {
			result[i] = data[start]
		}
	}
	return result
}

func RandomSampling(data []float64, targetLen int) []float64 {
	n := len(data)
	if targetLen >= n || targetLen <= 0 {
		return data
	}

	rand.Seed(time.Now().UnixNano())
	indices := make(map[int]bool)
	for len(indices) < targetLen {
		indices[rand.Intn(n)] = true
	}

	sortedIndices := make([]int, 0, targetLen)
	for idx := range indices {
		sortedIndices = append(sortedIndices, idx)
	}
	sort.Ints(sortedIndices)

	result := make([]float64, targetLen)
	for i, idx := range sortedIndices {
		result[i] = data[idx]
	}
	return result
}

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

// --- IO Utilities ---

func readTimeSeries(path string) ([]float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var reader *bufio.Scanner
	if strings.HasSuffix(path, ".gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		defer gr.Close()
		reader = bufio.NewScanner(gr)
	} else {
		reader = bufio.NewScanner(f)
	}

	var data []float64
	for reader.Scan() {
		line := strings.TrimSpace(reader.Text())
		if line == "" {
			continue
		}
		val, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return nil, err
		}
		data = append(data, val)
	}
	return data, reader.Err()
}

// --- Helper Functions ---

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func abs(a float64) float64 { return math.Abs(a) }

type Stats struct{ Mean, StdDev float64 }

func computeStats(data []float64) Stats {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	mean := sum / float64(len(data))
	sqDiff := 0.0
	for _, v := range data {
		sqDiff += (v - mean) * (v - mean)
	}
	return Stats{Mean: mean, StdDev: math.Sqrt(sqDiff / float64(len(data)))}
}

func zNormalize(data []float64) []float64 {
	stats := computeStats(data)
	if stats.StdDev == 0 {
		return append([]float64(nil), data...)
	}
	res := make([]float64, len(data))
	for i, v := range data {
		res[i] = (v - stats.Mean) / stats.StdDev
	}
	return res
}

// --- Distance Implementations ---

func EuclideanDistance(a, b []float64) float64 {
	sum := 0.0
	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

func ManhattanDistance(a, b []float64) float64 {
	sum := 0.0
	for i := 0; i < len(a); i++ {
		sum += abs(a[i] - b[i])
	}
	return sum
}

func ChebyshevDistance(a, b []float64) float64 {
	maxDiff := 0.0
	for i := 0; i < len(a); i++ {
		diff := abs(a[i] - b[i])
		if diff > maxDiff {
			maxDiff = diff
		}
	}
	return maxDiff
}

func DTW(a, b []float64, window int) float64 {
	n, m := len(a), len(b)
	prev := make([]float64, m+1)
	curr := make([]float64, m+1)

	for j := 0; j <= m; j++ {
		prev[j] = math.Inf(1)
	}
	prev[0] = 0

	for i := 1; i <= n; i++ {
		curr[0] = math.Inf(1)
		for j := 1; j <= m; j++ {
			if abs(float64(i-j)) > float64(window) {
				curr[j] = math.Inf(1)
				continue
			}
			cost := abs(a[i-1] - b[j-1])
			minPrev := prev[j]
			if prev[j-1] < minPrev {
				minPrev = prev[j-1]
			}
			if curr[j-1] < minPrev {
				minPrev = curr[j-1]
			}
			curr[j] = cost + minPrev
		}
		copy(prev, curr)
	}
	return prev[m]
}

func LCSS(a, b []float64, delta int, epsilon float64) float64 {
	n, m := len(a), len(b)
	prev := make([]int, m+1)
	curr := make([]int, m+1)

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if abs(float64(i-j)) <= float64(delta) && abs(a[i-1]-b[j-1]) <= epsilon {
				curr[j] = prev[j-1] + 1
			} else {
				curr[j] = max(prev[j], curr[j-1])
			}
		}
		copy(prev, curr)
	}
	minLen := float64(min(n, m))
	if minLen == 0 {
		return 1.0
	}
	return 1.0 - (float64(prev[m]) / minLen)
}

// --- MPdist Implementation ---

func MPdist(a, b []float64, window int, normalizeSubsequences bool) float64 {
	if len(a) < window || len(b) < window {
		return math.Inf(1)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var abProfile, baProfile []float64

	go func() { defer wg.Done(); abProfile = computeMatrixProfile(a, b, window, normalizeSubsequences) }()
	go func() { defer wg.Done(); baProfile = computeMatrixProfile(b, a, window, normalizeSubsequences) }()

	wg.Wait()
	jointProfile := append(abProfile, baProfile...)
	sort.Float64s(jointProfile)

	k := int(math.Ceil(0.05 * float64(len(jointProfile))))
	if k >= len(jointProfile) {
		k = len(jointProfile) - 1
	}
	if k < 0 {
		return 0.0
	}
	return jointProfile[k]
}

func computeMatrixProfile(query, target []float64, m int, zNorm bool) []float64 {
	qLen := len(query) - m + 1
	tLen := len(target) - m + 1
	result := make([]float64, qLen)

	numWorkers := 8
	workChan := make(chan int, qLen)
	var wg sync.WaitGroup

	for i := 0; i < qLen; i++ {
		workChan <- i
	}
	close(workChan)

	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			for i := range workChan {
				subQuery := query[i : i+m]
				minDist := math.Inf(1)

				var qMean, qStd float64
				if zNorm {
					stats := computeStats(subQuery)
					qMean, qStd = stats.Mean, stats.StdDev
				}

				for j := 0; j < tLen; j++ {
					subTarget := target[j : j+m]
					dist := 0.0

					if zNorm {
						tStats := computeStats(subTarget)
						if qStd == 0 || tStats.StdDev == 0 {
							dist = math.Sqrt(float64(m))
						} else {
							dotProd := 0.0
							for k := 0; k < m; k++ {
								dotProd += (subQuery[k] - qMean) * (subTarget[k] - tStats.Mean)
							}
							corr := (dotProd / float64(m)) / (qStd * tStats.StdDev)
							if corr > 1.0 {
								corr = 1.0
							}
							if corr < -1.0 {
								corr = -1.0
							}
							dist = math.Sqrt(2 * float64(m) * (1 - corr))
						}
					} else {
						for k := 0; k < m; k++ {
							d := subQuery[k] - subTarget[k]
							dist += d * d
						}
						dist = math.Sqrt(dist)
					}
					if dist < minDist {
						minDist = dist
					}
				}
				result[i] = minDist
			}
		}()
	}
	wg.Wait()
	return result
}
