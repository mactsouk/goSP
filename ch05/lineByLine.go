package main

import (
	"bufio"
	"fmt"
	"os"
)

// readLines reads a text file line by line and returns a slice of lines or an error
func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func main() {
	var input string
	if len(os.Args) > 1 {
		input = os.Args[1]
	}

	lines, err := readLines(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	for i, line := range lines {
		fmt.Printf("%d: %s\n", i+1, line)
	}
}
