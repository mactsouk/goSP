// Package tslib provides a collection of algorithms for time series analysis,
// including distance metrics (Euclidean, DTW, MPdist) and data summarization techniques.
package tslib

import (
	"math/rand"
	"sort"
	"time"
)

// PAA (Piecewise Aggregate Approximation) reduces the dimensionality of a time series
// by dividing it into 'targetLen' equal-sized segments and replacing the data
// in each segment with its mean value.
//
// This method reduces noise and storage requirements while preserving the overall
// trend of the data. If targetLen is greater than or equal to the length of the
// input data, the original slice is returned.
//
// Parameters:
//   - data: The input time series slice.
//   - targetLen: The desired length of the summarized series.
//
// Returns a new slice of float64 containing the reduced time series.
func PAA(data []float64, targetLen int) []float64 {
	n := len(data)
	if targetLen >= n || targetLen <= 0 {
		return data
	}

	result := make([]float64, targetLen)
	// Calculate segment size as a float to handle non-integer divisions evenly
	segmentSize := float64(n) / float64(targetLen)

	for i := 0; i < targetLen; i++ {
		// Determine the start and end indices for the current segment
		start := int(float64(i) * segmentSize)
		end := int(float64(i+1) * segmentSize)

		// Clamp the end index to prevent out-of-bounds errors on the final segment
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
			// Fallback for empty segments (rare edge case with extreme reduction)
			result[i] = data[start]
		}
	}
	return result
}

// RandomSampling reduces a time series by selecting 'targetLen' data points
// at random. Crucially, this function sorts the selected indices before retrieving
// values, ensuring the chronological order of the time series is preserved.
//
// This technique is useful for creating a sparse representation of the signal
// without introducing the smoothing artifacts of PAA.
//
// Parameters:
//   - data: The input time series slice.
//   - targetLen: The number of points to select.
//
// Returns a new slice containing the sampled values in their original relative order.
func RandomSampling(data []float64, targetLen int) []float64 {
	n := len(data)
	if targetLen >= n || targetLen <= 0 {
		return data
	}

	// Seed the random number generator to ensure varied results across runs
	rand.Seed(time.Now().UnixNano())

	// Use a map to ensure unique indices are selected
	indices := make(map[int]bool)
	for len(indices) < targetLen {
		indices[rand.Intn(n)] = true
	}

	// Extract keys and sort them to restore temporal ordering
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

// Stepping (also known as Decimation) reduces the time series by keeping
// only every 'step'-th data point. For example, a step of 2 keeps every
// second point.
//
// This is the fastest summarization method but it is prone to
// "aliasing", where high-frequency features (like sharp spikes)
// falling between steps are completely lost.
//
// Parameters:
//   - data: The input time series slice.
//   - step: The interval at which to sample points (must be > 1).
//
// Returns a new slice containing the decimated time series.
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
