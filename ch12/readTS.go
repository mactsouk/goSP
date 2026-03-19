package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"strconv"
)

func readGzipTimeSeries(path string) ([]float64, error) {
	// 1. Open the file on disk
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 2. Wrap the file reader in a gzip reader
	gr, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	// Important: Close the gzip reader to flush any buffers
	defer gr.Close()

	var values []float64

	// 3. Scan from the gzip reader (gr), not the raw file
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <file.gz>")
		os.Exit(1)
	}

	filename := os.Args[1]

	data, err := readGzipTimeSeries(filename)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	fmt.Printf("Successfully read %d values from %s\n", len(data), filename)
	if len(data) > 0 {
		fmt.Printf("First value: %f\n", data[0])
		fmt.Printf("Last value:  %f\n", data[len(data)-1])
	}
}
