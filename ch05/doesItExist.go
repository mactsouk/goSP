package main

import (
	"fmt"
	"os"
    "errors"
    "io/fs"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<path>")
		os.Exit(1)
	}
	path := os.Args[1]

	_, err := os.Stat(path)
	if err == nil {
		fmt.Printf("Path %s exists.\n", path)
		return
	}

    // Thanks Natan!
	if errors.Is(err, fs.ErrNotExist) {
		fmt.Printf("Path %s does not exist.\n", path)
		return
	}

	fmt.Fprintf(os.Stderr, "Error checking path: %v\n", err)
	os.Exit(1)
}
