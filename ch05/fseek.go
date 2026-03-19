package main

import (
	"flag"
	"fmt"
	"os"
    "io"
	"strconv"
)

// seekAndRead seeks to the given offset in the file and reads length bytes.
func seekAndRead(fileName string, offset int64, length int) ([]byte, error) {
	// Open file
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	// Seek to offset
	_, err = f.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("error seeking file: %w", err)
	}

	// Read requested bytes
	buf := make([]byte, length)
	n, err := f.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	return buf[:n], nil
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 3 {
		fmt.Println("Usage: seekread <file> <offset> <length>")
		os.Exit(1)
	}
	fileName := args[0]

	offset, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid offset: %v\n", err)
		os.Exit(1)
	}

	length, err := strconv.Atoi(args[2])
	if err != nil || length <= 0 {
		fmt.Fprintf(os.Stderr, "Invalid length: %v\n", args[2])
		os.Exit(1)
	}

	data, err := seekAndRead(fileName, offset, length)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Read %d bytes from offset %d:\n", len(data), offset)
	fmt.Println(string(data))
}
