package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyResult represents the result of copying to a single directory
type CopyResult struct {
	Directory string
	Success   bool
	Error     string
}

// copyFileToDirectories copies a source file to multiple target directories
// Returns a slice of CopyResult showing the outcome for each directory
func copyFileToDirectories(sourceFile string, targetDirs []string) []CopyResult {
	results := make([]CopyResult, len(targetDirs))

	// Check if source file exists and is readable
	sourceInfo, err := os.Stat(sourceFile)
	if err != nil {
		// If source doesn't exist, mark all operations as failed
		for i, dir := range targetDirs {
			results[i] = CopyResult{
				Directory: dir,
				Success:   false,
				Error:     fmt.Sprintf("source file error: %v", err),
			}
		}
		return results
	}

	if sourceInfo.IsDir() {
		// If source is a directory, mark all operations as failed
		for i, dir := range targetDirs {
			results[i] = CopyResult{
				Directory: dir,
				Success:   false,
				Error:     "source is a directory, not a file",
			}
		}
		return results
	}

	// Extract the filename from the source path
	filename := filepath.Base(sourceFile)

	// Process each target directory
	for i, targetDir := range targetDirs {
		results[i] = copyToSingleDirectory(sourceFile, targetDir, filename)
	}

	return results
}

// copyToSingleDirectory handles copying to one specific directory
func copyToSingleDirectory(sourceFile, targetDir, filename string) CopyResult {
	result := CopyResult{Directory: targetDir}

	// Check if target directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		result.Success = false
		result.Error = "directory does not exist"
		return result
	}

	// Construct target file path
	targetFile := filepath.Join(targetDir, filename)

	// Check if target file already exists
	if _, err := os.Stat(targetFile); err == nil {
		result.Success = false
		result.Error = "file already exists in directory"
		return result
	}

	// Perform the actual copy operation
	err := copyFile(sourceFile, targetFile)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("copy failed: %v", err)
		return result
	}

	result.Success = true
	return result
}

// copyFile performs the actual file copying operation
func copyFile(src, dst string) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}
	defer destFile.Close()

	// Copy the content
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy content: %w", err)
	}

	// Sync to ensure data is written to disk
	err = destFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination: %w", err)
	}

	return nil
}

func main() {
	// Check command line arguments
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <source_file> <target_dir1> [target_dir2] ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nCopies source_file to multiple target directories.\n")
		fmt.Fprintf(os.Stderr, "Continues processing even if some operations fail.\n")
		fmt.Fprintf(os.Stderr, "Example:\n")
		fmt.Fprintf(os.Stderr, "  %s config.txt /etc /home/user/backup /tmp/configs\n", os.Args[0])
		os.Exit(1)
	}

	sourceFile := os.Args[1]
	targetDirs := os.Args[2:]
	fmt.Printf("Copying '%s' to %d directories..\n", sourceFile, len(targetDirs))
	results := copyFileToDirectories(sourceFile, targetDirs)

	successCount := 0
	failureCount := 0

	fmt.Println("Results:")
	fmt.Println("--------")

	for _, result := range results {
		if result.Success {
			fmt.Printf("%s - SUCCESS\n", result.Directory)
			successCount++
		} else {
			fmt.Printf("%s - FAILED: %s\n", result.Directory, result.Error)
			failureCount++
		}
	}

	fmt.Printf("Summary:\n")
	fmt.Printf("--------\n")
	fmt.Printf("Total directories: %d\n", len(results))
	fmt.Printf("Successful copies: %d\n", successCount)
	fmt.Printf("Failed copies: %d\n", failureCount)

	if failureCount > 0 {
		fmt.Printf("Note: Some operations failed. Check directory existence and permissions.\n")
		if successCount == 0 {
			os.Exit(1) // All failed
		} else {
			os.Exit(2) // Partial success
		}
	}

	fmt.Printf("All copy operations completed successfully!\n")
}
