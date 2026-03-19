package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func which(command string) (string, error) {
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, dir := range paths {
		fullPath := filepath.Join(dir, command)
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}

		// Check if the file is executable
		if fileInfo.Mode().Perm()&0111 != 0 {
			return fullPath, nil
		}
	}

	return "", fmt.Errorf("%s not found in PATH", command)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command>\n", os.Args[0])
		os.Exit(1)
	}

	command := os.Args[1]
	path, err := which(command)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(path)
}
