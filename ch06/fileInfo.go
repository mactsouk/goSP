// fileinfo.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(filepath.Base(os.Args[0]), "<file-path>")
		os.Exit(1)
	}
	path := os.Args[1]

	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File:", path)
	fmt.Println("Name:", info.Name())
	fmt.Println("Size:", info.Size(), "bytes")
	fmt.Println("Permissions:", info.Mode())
	fmt.Println("Modified:", info.ModTime().Format(time.RFC1123))
	if info.IsDir() {
		fmt.Println("Type: Directory")
	} else {
		fmt.Println("Type: File")
	}
}
