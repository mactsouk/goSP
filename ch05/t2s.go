package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseMode(mode string) (direction string, count int, err error) {
	if strings.HasSuffix(mode, "tabs") {
		numStr := strings.TrimSuffix(mode, "tabs")
		count, err = strconv.Atoi(numStr)
		if err != nil {
			return "", 0, fmt.Errorf("invalid mode: %s", mode)
		}
		direction = "spaces-to-tabs"
		return
	} else if strings.HasPrefix(mode, "tabs") {
		numStr := strings.TrimPrefix(mode, "tabs")
		count, err = strconv.Atoi(numStr)
		if err != nil {
			return "", 0, fmt.Errorf("invalid mode: %s", mode)
		}
		direction = "tabs-to-spaces"
		return
	}
	return "", 0, fmt.Errorf("unknown mode: %s", mode)
}

func main() {
	// Mode flag like "2tabs" or "tabs4"
	mode := flag.String("mode", "tabs4", "2tabs | tabs4")
	flag.Parse()

	// Positional arguments: input and optional output file
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: converter -mode=2tabs input.txt [output.txt]")
		os.Exit(1)
	}

	inputFile := args[0]
	var outputFile string
	if len(args) > 1 {
		outputFile = args[1]
	}

	// Parse mode
	direction, count, err := parseMode(*mode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Open input
	in, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer in.Close()

	// Prepare output
	var out *os.File
	if outputFile == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer out.Close()
	}

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	spaceStr := strings.Repeat(" ", count)

	for scanner.Scan() {
		line := scanner.Text()
		var converted string
		switch direction {
		case "tabs-to-spaces":
			converted = strings.ReplaceAll(line, "\t", spaceStr)
		case "spaces-to-tabs":
			converted = strings.ReplaceAll(line, spaceStr, "\t")
		}
		fmt.Fprintln(writer, converted)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
}
