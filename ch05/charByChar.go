package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// readChars reads a text input character by character and returns a slice of runes
func readChars(file *os.File) ([]rune, error) {
	reader := bufio.NewReader(file)
	var chars []rune

	for {
		r, _, err := reader.ReadRune() // Read one UTF-8 character
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		chars = append(chars, r)
	}

	return chars, nil
}

func main() {
	var input *os.File
	var err error

	if len(os.Args) > 1 {
		input, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer input.Close()
	} else {
		input = os.Stdin
	}

	chars, err := readChars(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	for i, c := range chars {
		fmt.Printf("%d: %c\n", i+1, c)
	}
}
