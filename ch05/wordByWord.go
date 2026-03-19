package main

import (
	"bufio"
	"fmt"
	"os"
)

// readWords reads a text file word by word and returns a slice of words or an error
func readWords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords) // Split by words

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func main() {
	var input string
	if len(os.Args) > 1 {
		input = os.Args[1]
	}

	words, err := readWords(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	for i, word := range words {
		fmt.Printf("%d: %s\n", i+1, word)
	}
}
