package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

var (
	countLines = flag.Bool("l", false, "print the newline counts")
	countWords = flag.Bool("w", false, "print the word counts")
	countBytes = flag.Bool("c", false, "print the byte counts")
	countRunes = flag.Bool("m", false, "print the character (rune) counts")
)

func countWordsInLine(s string) int {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}

func count(r io.Reader) (lines, words, bytes, runes int) {
	buf := bufio.NewReader(r)

	for {
		line, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "read error: %v\n", err)
			break
		}

		if len(line) > 0 {
			lines++
			trimmed := strings.TrimRight(line, "\r\n")
			words += countWordsInLine(trimmed)
			bytes += len(line)
			runes += utf8.RuneCountInString(trimmed) + 1 // include newline rune
		}

		if err == io.EOF {
			break
		}
	}
	return
}

func printCounts(lines, words, bytes, runes int, filename string) {
	// If no flags are set, default to all (l, w, c)
	if !*countLines && !*countWords && !*countBytes && !*countRunes {
		fmt.Printf("%7d %7d %7d %s\n", lines, words, bytes, filename)
		return
	}

	if *countLines {
		fmt.Printf("%7d ", lines)
	}
	if *countWords {
		fmt.Printf("%7d ", words)
	}
	if *countBytes {
		fmt.Printf("%7d ", bytes)
	}
	if *countRunes {
		fmt.Printf("%7d ", runes)
	}
	fmt.Printf("%s\n", filename)
}

func main() {
	flag.Parse()
	files := flag.Args()

	if len(files) == 0 {
		lines, words, bytes, runes := count(os.Stdin)
		printCounts(lines, words, bytes, runes, "")
		return
	}

	var totalLines, totalWords, totalBytes, totalRunes int

	for _, filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "wc: %s: %v\n", filename, err)
			continue
		}
		defer file.Close()

		lines, words, bytes, runes := count(file)
		printCounts(lines, words, bytes, runes, filename)

		totalLines += lines
		totalWords += words
		totalBytes += bytes
		totalRunes += runes
	}

	if len(files) > 1 {
		printCounts(totalLines, totalWords, totalBytes, totalRunes, "total")
	}
}
