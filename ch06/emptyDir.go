// find_empty_dirs_fs.go
package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", filepath.Base(os.Args[0]), "<directory>")
		os.Exit(1)
	}
	root := os.Args[1]

	// Walk the directory tree
	err := fs.WalkDir(os.DirFS(root), ".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", path, err)
				return nil
			}

			if d.IsDir() {
				entries, err := os.ReadDir(filepath.Join(root, path))
				if err != nil {
					fmt.Fprintf(os.Stderr,
						"Error in directory %s: %v\n", path, err)
					return nil
				}
				if len(entries) == 0 {
					fmt.Println(filepath.Join(root, path))
				}
			}
			return nil
		})

	if err != nil {
		log.Fatalf("Error walking directory tree: %v\n", err)
	}
}
