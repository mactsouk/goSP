package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"
)

type Statistics struct {
	Count    int
	Min      float64
	Max      float64
	Mean     float64
	Median   float64
	Q1       float64
	Q3       float64
	StdDev   float64
	Skewness float64
	Kurtosis float64
}

func readGzipTimeSeries(path string) ([]float64, error) {
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

func percentile(sortedData []float64, p float64) float64 {
	if len(sortedData) == 0 {
		return 0
	}
	idx := p * float64(len(sortedData)-1)
	lower := int(math.Floor(idx))
	upper := int(math.Ceil(idx))

	if lower == upper {
		return sortedData[lower]
	}

	fraction := idx - float64(lower)
	return sortedData[lower] + (fraction * (sortedData[upper] - sortedData[lower]))
}

func computeStats(data []float64) Statistics {
	if len(data) == 0 {
		return Statistics{}
	}

	n := float64(len(data))
	minVal := data[0]
	maxVal := data[0]
	sum := 0.0

	// Pass 1: Min, Max, Sum
	for _, v := range data {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
		sum += v
	}

	mean := sum / n

	// Pass 2: Moments (Variance, Skewness, Kurtosis)
	var sumSqDiff, sumCubedDiff, sumQuartDiff float64

	for _, v := range data {
		diff := v - mean
		diffSq := diff * diff

		sumSqDiff += diffSq
		sumCubedDiff += diffSq * diff
		sumQuartDiff += diffSq * diffSq
	}

	// Population Variance and StdDev
	variance := sumSqDiff / n
	stdDev := math.Sqrt(variance)

	var skewness, kurtosis float64

	// Avoid division by zero if dataset is constant (stdDev = 0)
	if stdDev > 0 {
		// Skewness = (Mean of cubed diffs) / (StdDev^3)
		skewness = (sumCubedDiff / n) / math.Pow(stdDev, 3)

		// Excess Kurtosis = (Mean of quart diffs) / (StdDev^4) - 3
		kurtosis = ((sumQuartDiff / n) / math.Pow(stdDev, 4)) - 3.0
	}

	// Sort copy for percentiles
	sortedData := make([]float64, len(data))
	copy(sortedData, data)
	sort.Float64s(sortedData)

	return Statistics{
		Count:    len(data),
		Min:      minVal,
		Max:      maxVal,
		Mean:     mean,
		Median:   percentile(sortedData, 0.5),
		Q1:       percentile(sortedData, 0.25),
		Q3:       percentile(sortedData, 0.75),
		StdDev:   stdDev,
		Skewness: skewness,
		Kurtosis: kurtosis,
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run stats.go <data.gz>")
		os.Exit(1)
	}

	filename := os.Args[1]
	data, err := readGzipTimeSeries(filename)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	if len(data) == 0 {
		fmt.Println("File was empty.")
		return
	}

	stats := computeStats(data)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Metric\tValue")
	fmt.Fprintln(w, "------\t-----")
	fmt.Fprintf(w, "Count\t%d\n", stats.Count)
	fmt.Fprintf(w, "Min\t%.4f\n", stats.Min)
	fmt.Fprintf(w, "Q1 (25%%)\t%.4f\n", stats.Q1)
	fmt.Fprintf(w, "Median\t%.4f\n", stats.Median)
	fmt.Fprintf(w, "Q3 (75%%)\t%.4f\n", stats.Q3)
	fmt.Fprintf(w, "Max\t%.4f\n", stats.Max)
	fmt.Fprintf(w, "Mean\t%.4f\n", stats.Mean)
	fmt.Fprintf(w, "StdDev\t%.4f\n", stats.StdDev)
	fmt.Fprintf(w, "Skewness\t%.4f\n", stats.Skewness)
	fmt.Fprintf(w, "Kurtosis\t%.4f\n", stats.Kurtosis)
	w.Flush()
}
