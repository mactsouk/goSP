package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

type Stats struct {
	Dirs  int
	Files int
	Bytes int64
}

type FileInfo struct {
	Path string
	Size int64
}

func main() {
	// Command-line flags
	topN := flag.Int("top", 5, "Number of largest files to display")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <directory>\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(1)
	}

	root := flag.Arg(0)
	var stats Stats
	var files []FileInfo

	err := fs.WalkDir(os.DirFS(root), ".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", path, err)
				return nil
			}

			if d.IsDir() {
				stats.Dirs++
				return nil
			}

			info, err := d.Info()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting info for %s: %v\n", path, err)
				return nil
			}

			stats.Files++
			stats.Bytes += info.Size()
			files = append(files,
				FileInfo{Path: filepath.Join(root, path), Size: info.Size()})
			return nil
		})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
		os.Exit(1)
	}

	// Sort files by size descending
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size > files[j].Size
	})

	fmt.Printf("Directory statistics for %s:\n", root)
	fmt.Printf("Directories: %d\n", stats.Dirs)
	fmt.Printf("Files: %d\n", stats.Files)
	fmt.Printf("Total size: %d bytes\n\n", stats.Bytes)

	fmt.Printf("Top %d largest files:\n", *topN)
	for i, f := range files {
		if i >= *topN {
			break
		}
		fmt.Printf("%d. %s (%d bytes)\n", i+1, f.Path, f.Size)
	}
}
