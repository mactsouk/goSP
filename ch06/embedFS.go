// embedded_explorer.go
package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

//go:embed assets/*
var embeddedFiles embed.FS

func main() {
	// Command-line flag for optional filename to read
	fileToRead := flag.String("file", "", "Read and display contents of a specific embedded file")
	flag.Parse()

	// Create a sub-filesystem rooted at "assets"
	assetsFS, err := fs.Sub(embeddedFiles, "assets")
	if err != nil {
		log.Fatalf("failed to get sub FS: %v", err)
	}

	// Walk the embedded filesystem
	err = fs.WalkDir(assetsFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", path, err)
			return nil
		}

		info, _ := d.Info()
		if d.IsDir() {
			fmt.Printf("[DIR]  %s\n", path)
		} else {
			fmt.Printf("[FILE] %s (%d bytes)\n", path, info.Size())
		}
		return nil
	})
	if err != nil {
		log.Fatalf("walk error: %v", err)
	}

	// If a filename was provided, read and display it
	if *fileToRead != "" {
		content, err := fs.ReadFile(assetsFS, filepath.Clean(*fileToRead))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", *fileToRead, err)
			os.Exit(1)
		}
		fmt.Printf("\nContents of %s:\n", *fileToRead)
		fmt.Println(string(content))
	}
}
