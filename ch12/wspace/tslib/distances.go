package tslib

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"sync"
)

// Internal helper for min integer
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// --- Distance Functions ---
func EuclideanDistance(a, b []float64) float64 {
	n := min(len(a), len(b))
	sum := 0.0
	for i := 0; i < n; i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

func ManhattanDistance(a, b []float64) float64 {
	n := min(len(a), len(b))
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += math.Abs(a[i] - b[i])
	}
	return sum
}

func ChebyshevDistance(a, b []float64) float64 {
	n := min(len(a), len(b))
	maxDiff := 0.0
	for i := 0; i < n; i++ {
		diff := math.Abs(a[i] - b[i])
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
			if math.Abs(float64(i-j)) > float64(window) {
				curr[j] = math.Inf(1)
				continue
			}
			cost := math.Abs(a[i-1] - b[j-1])
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
			if math.Abs(float64(i-j)) <= float64(delta) && math.Abs(a[i-1]-b[j-1]) <= epsilon {
				curr[j] = prev[j-1] + 1
			} else {
				if prev[j] > curr[j-1] {
					curr[j] = prev[j]
				} else {
					curr[j] = curr[j-1]
				}
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

func MPdist(a, b []float64, window int) float64 {
	if len(a) < window || len(b) < window {
		return math.Inf(1)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var abProfile, baProfile []float64

	go func() {
		defer wg.Done()
		abProfile = computeMatrixProfile(a, b, window)
	}()

	go func() {
		defer wg.Done()
		baProfile = computeMatrixProfile(b, a, window)
	}()

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

func computeMatrixProfile(query, target []float64, m int) []float64 {
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

				for j := 0; j < tLen; j++ {
					subTarget := target[j : j+m]
					dist := 0.0
					for k := 0; k < m; k++ {
						d := subQuery[k] - subTarget[k]
						dist += d * d
					}
					dist = math.Sqrt(dist)
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

// --- IO Utility Functions ---

// ReadGzipTimeSeries reads float64 values from a .gz file.
// It is exported so it can be used by the main application.
func ReadGzipTimeSeries(path string) ([]float64, error) {
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
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return values, nil
}
