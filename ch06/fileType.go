// filetype.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func fileType(info os.FileInfo) string {
	mode := info.Mode()

	switch {
	case mode.IsRegular():
		return "regular file"
	case mode.IsDir():
		return "directory"
	case mode&os.ModeSymlink != 0:
		return "symbolic link"
	case mode&os.ModeNamedPipe != 0:
		return "named pipe (FIFO)"
	case mode&os.ModeSocket != 0:
		return "socket"
	case mode&os.ModeDevice != 0:
		if mode&os.ModeCharDevice != 0 {
			return "character device"
		}
		return "block device"
	default:
		return "unknown"
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%s <file-path>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	path := os.Args[1]

	info, err := os.Lstat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Path: %s\n", path)
	fmt.Printf("Type: %s\n", fileType(info))
}
