package main

import (
	"fmt"
	"os"
	"strings"
)

func basename(path string, suffixes ...string) string {
	// Handle root directory case: "/" → "/"
	if path == "/" {
		return "/"
	}

	// Remove trailing slashes (except if the path is just "/")
	path = strings.TrimRight(path, "/")
	if path == "" { // Happens if input was something like "////"
		return "/"
	}

	// Extract the last component after the last '/'
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash >= 0 {
		path = path[lastSlash+1:]
	}

	// Handle empty result (shouldn't happen due to earlier checks)
	if path == "" {
		return "/"
	}

	// Remove suffix if specified (exact match)
	if len(suffixes) > 0 {
		for _, suffix := range suffixes {
			if strings.HasSuffix(path, suffix) {
				path = strings.TrimSuffix(path, suffix)
				break
			}
		}
	}

	return path
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path> [suffix]\n", os.Args[0])
		os.Exit(1)
	}

	path := os.Args[1]
	suffix := ""
	if len(os.Args) > 2 {
		suffix = os.Args[2]
	}

	result := basename(path, suffix)
	fmt.Println(result)
}
