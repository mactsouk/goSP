package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

// ImageProcessor simulates processing images with workers
type ImageProcessor struct {
	results chan string
}

func (p *ImageProcessor) Process(images []string) ([]string, error) {
	p.results = make(chan string) // Unbuffered - the leak source!

	// Spawn workers to process each image
	for _, img := range images {
		go func(filename string) {
			// Simulate image processing (resizing, filtering, etc.)
			time.Sleep(20 * time.Millisecond)

			// Simulate occasional corruption
			if filename == "corrupted.jpg" {
				p.results <- "ERROR: " + filename
				return
			}

			p.results <- "✓ Processed: " + filename
		}(img)
	}

	// Collect results
	var processed []string
	for range images {
		result := <-p.results
		if result[:5] == "ERROR" {
			// Oh no! Early return without draining remaining results
			// Workers still trying to send will block forever
			return nil, fmt.Errorf("image processing failed: %s", result)
		}
		processed = append(processed, result)
	}
	return processed, nil
}

func main() {
	processor := &ImageProcessor{}

	images := []string{
		"photo1.jpg",
		"photo2.jpg",
		"corrupted.jpg", // This will cause an error
		"photo3.jpg",    // These workers will leak!
		"photo4.jpg",
	}

	fmt.Println("Processing images...")
	results, err := processor.Process(images)

	if err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println("* * * Some worker goroutines are now leaked!\n")
	} else {
		fmt.Printf("✓ Processed %d images\n", len(results))
	}

	// Give GC time to detect leaks
	time.Sleep(100 * time.Millisecond)

	// Check for leaks
	fmt.Println("--- Goroutine Leak Profile ---")
	prof := pprof.Lookup("goroutineleak")

	if prof == nil {
		fmt.Println("\nProfile not available!")
		fmt.Println("Build with: GOEXPERIMENT=goroutineleakprofile go run main.go")
		os.Exit(1)
	}

	// Count leaks
	var buf []byte
	if err := prof.WriteTo(&writerFunc{fn: func(p []byte) (int, error) {
		buf = append(buf, p...)
		return len(p), nil
	}}, 1); err == nil {
		if len(buf) > 0 {
			fmt.Printf("\nDetected leaked goroutines!\n\n")
		}
	}

	// Show detailed traces
	if err := prof.WriteTo(os.Stdout, 2); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// fmt.Println("\nThe Fix:")
	// fmt.Println("   Use buffered channel: make(chan string, len(images))")
	// fmt.Println("   Or drain remaining results before returning")
}

// Helper to capture profile output
type writerFunc struct {
	fn func([]byte) (int, error)
}

func (w *writerFunc) Write(p []byte) (int, error) {
	return w.fn(p)
}
