package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

var (
	filename     = flag.String("file", "", "File to read for profiling")
	cpuProfile   = flag.String("cpuprofile", "", "Write CPU profile to file")
	memProfile   = flag.String("memprofile", "", "Write memory profile to file")
	method       = flag.String("method", "both", "Profiling method: unbuffered, buffered, or both")
	iterations   = flag.Int("iterations", 1, "Number of iterations to run")
	generateFile = flag.String("generate", "", "Generate sample file: small, medium, large, or custom size (e.g., 10MB)")
)

func main() {
	flag.Parse()

	// Handle sample file generation
	if *generateFile != "" {
		generateSampleFile(*generateFile)
		return
	}

	if *filename == "" {
		fmt.Println("Usage: go run <source.go> -file=<filename> [options]")
		fmt.Println("   or: go run <source.go> -generate=<size>")
		fmt.Println("\nOptions:")
		fmt.Println("  -file=<filename>        File to profile (required unless generating)")
		fmt.Println("  -generate=<size>        Generate sample file: small, medium, large, or size (e.g., 10MB)")
		fmt.Println("  -cpuprofile=<file>      Write CPU profile to file")
		fmt.Println("  -memprofile=<file>      Write memory profile to file")
		fmt.Println("  -method=<method>        unbuffered, buffered, or both (default: both)")
		fmt.Println("  -iterations=<n>         Number of iterations (default: 1)")
		fmt.Println("\nExample usage:")
		fmt.Println("  go run <source.go> -generate=large")
		fmt.Println("  go run <source.go> -file=sample_large.txt -cpuprofile=cpu.prof -memprofile=mem.prof")
		fmt.Println("  go tool pprof cpu.prof")
		fmt.Println("  go tool pprof mem.prof")
		os.Exit(1)
	}

	// Verify file exists
	fileInfo, err := os.Stat(*filename)
	if err != nil {
		log.Fatalf("Error accessing file: %v", err)
	}

	fmt.Printf("=== Go pprof I/O Profiler ===\n")
	fmt.Printf("File: %s\n", *filename)
	fmt.Printf("File size: %s (%d bytes)\n", formatBytes(fileInfo.Size()), fileInfo.Size())
	fmt.Printf("Method: %s\n", *method)
	fmt.Printf("Iterations: %d\n", *iterations)
	fmt.Printf("Go version: %s\n", runtime.Version())

	// Start CPU profiling if requested
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
		fmt.Printf("CPU profiling enabled: %s\n", *cpuProfile)
	}

	fmt.Println("\nStarting profiling...")

	// Run the specified profiling method
	switch *method {
	case "unbuffered":
		runUnbufferedProfiling(*filename, *iterations)
	case "buffered":
		runBufferedProfiling(*filename, *iterations)
	case "both":
		runBothProfiling(*filename, *iterations)
	default:
		log.Fatalf("Invalid method: %s", *method)
	}

	// Write memory profile if requested
	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		runtime.GC() // Force garbage collection for accurate memory profile
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Memory profile written: %s\n", *memProfile)
	}

	fmt.Println("\nProfiling complete!")
	if *cpuProfile != "" || *memProfile != "" {
		fmt.Println("\nAnalyze profiles with:")
		if *cpuProfile != "" {
			fmt.Printf("  go tool pprof %s\n", *cpuProfile)
		}
		if *memProfile != "" {
			fmt.Printf("  go tool pprof %s\n", *memProfile)
		}
		fmt.Println("\nUseful pprof commands:")
		fmt.Println("  (pprof) top     - Show top functions by CPU/memory usage")
		fmt.Println("  (pprof) list    - Show source code with annotations")
		fmt.Println("  (pprof) web     - Generate web-based visualization")
		fmt.Println("  (pprof) png     - Generate PNG visualization")
	}
}

func runUnbufferedProfiling(filename string, iterations int) {
	fmt.Printf("\n--- Unbuffered Reading Profile ---\n")
	start := time.Now()
	totalBytes := int64(0)

	for i := 0; i < iterations; i++ {
		bytes := readUnbuffered(filename)
		totalBytes += bytes
		if i == 0 {
			fmt.Printf("First iteration: %d bytes read\n", bytes)
		}
	}

	duration := time.Since(start)
	throughput := float64(totalBytes) / duration.Seconds() / (1024 * 1024)

	fmt.Printf("Total time: %v\n", duration)
	fmt.Printf("Total bytes: %s\n", formatBytes(totalBytes))
	fmt.Printf("Throughput: %.2f MB/s\n", throughput)
}

func runBufferedProfiling(filename string, iterations int) {
	fmt.Printf("\n--- Buffered Reading Profile ---\n")
	start := time.Now()
	totalBytes := int64(0)

	for i := 0; i < iterations; i++ {
		bytes := readBuffered(filename)
		totalBytes += bytes
		if i == 0 {
			fmt.Printf("First iteration: %d bytes read\n", bytes)
		}
	}

	duration := time.Since(start)
	throughput := float64(totalBytes) / duration.Seconds() / (1024 * 1024)

	fmt.Printf("Total time: %v\n", duration)
	fmt.Printf("Total bytes: %s\n", formatBytes(totalBytes))
	fmt.Printf("Throughput: %.2f MB/s\n", throughput)
}

func runBothProfiling(filename string, iterations int) {
	// Run unbuffered first
	fmt.Printf("\n--- Phase 1: Unbuffered Reading ---\n")
	unbufferedStart := time.Now()
	unbufferedBytes := int64(0)

	for i := 0; i < iterations; i++ {
		bytes := readUnbuffered(filename)
		unbufferedBytes += bytes
	}
	unbufferedDuration := time.Since(unbufferedStart)
	unbufferedThroughput := float64(unbufferedBytes) / unbufferedDuration.Seconds() / (1024 * 1024)

	// Force GC between tests for cleaner profiling
	runtime.GC()
	runtime.GC()

	// Run buffered second
	fmt.Printf("\n--- Phase 2: Buffered Reading ---\n")
	bufferedStart := time.Now()
	bufferedBytes := int64(0)

	for i := 0; i < iterations; i++ {
		bytes := readBuffered(filename)
		bufferedBytes += bytes
	}
	bufferedDuration := time.Since(bufferedStart)
	bufferedThroughput := float64(bufferedBytes) / bufferedDuration.Seconds() / (1024 * 1024)

	// Display comparison
	fmt.Printf("\n--- Comparison Results ---\n")
	fmt.Printf("Unbuffered: %v (%.2f MB/s)\n", unbufferedDuration, unbufferedThroughput)
	fmt.Printf("Buffered:   %v (%.2f MB/s)\n", bufferedDuration, bufferedThroughput)

	if bufferedDuration > 0 {
		speedup := float64(unbufferedDuration) / float64(bufferedDuration)
		fmt.Printf("Speedup:    %.2fx faster with buffered I/O\n", speedup)
	}
}

// readUnbuffered performs byte-by-byte unbuffered reading
func readUnbuffered(filename string) int64 {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	buffer := make([]byte, 1)
	totalBytes := int64(0)

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			totalBytes += int64(n)
			// Simulate some processing work
			_ = buffer[0]
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Read error: %v", err)
		}
	}

	return totalBytes
}

// readBuffered performs buffered reading
func readBuffered(filename string) int64 {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	totalBytes := int64(0)

	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Read error: %v", err)
		}
		totalBytes++
		// Simulate some processing work
		_ = b
	}

	return totalBytes
}

// readBufferedChunk performs chunk-based buffered reading (alternative method)
func readBufferedChunk(filename string) int64 {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, 4096) // 4KB chunks
	totalBytes := int64(0)

	for {
		n, err := reader.Read(buffer)
		if n > 0 {
			totalBytes += int64(n)
			// Simulate processing the chunk
			for i := 0; i < n; i++ {
				_ = buffer[i]
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Read error: %v", err)
		}
	}

	return totalBytes
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// generateSampleFile creates sample files for testing
func generateSampleFile(sizeSpec string) {
	var size int64
	var filename string

	switch sizeSpec {
	case "small":
		size = 1 * 1024 // 1KB
		filename = "/tmp/sample_small.txt"
	case "medium":
		size = 100 * 1024 // 100KB
		filename = "/tmp/sample_medium.txt"
	case "large":
		size = 10 * 1024 * 1024 // 10MB
		filename = "/tmp/sample_large.txt"
	default:
		// Try to parse custom size (e.g., "5MB", "500KB", "2GB")
		customSize, customFilename := parseCustomSize(sizeSpec)
		if customSize == 0 {
			log.Fatalf("Invalid size specification: %s", sizeSpec)
		}
		size = customSize
		filename = customFilename
	}

	fmt.Printf("Generating sample file: %s (%s)\n", filename, formatBytes(size))

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// Generate varied content to make reading more realistic
	patterns := []string{
		"The quick brown fox jumps over the lazy dog. ",
		"Go is an open source programming language that makes it easy to build simple, reliable, and efficient software. ",
		"Buffered I/O can significantly improve performance by reducing system calls. ",
		"Profiling helps identify bottlenecks and optimize application performance. ",
		"Package bufio implements buffered I/O. It wraps an io.Reader or io.Writer object. ",
		"The pprof package serves via its HTTP server runtime profiling data. ",
		"1234567890 !@#$%^&*()_+-=[]{}|;':\",./<>? abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ\n",
	}

	bytesWritten := int64(0)
	lineNumber := 1

	for bytesWritten < size {
		// Add line numbers every 10 lines for structure
		if lineNumber%10 == 1 {
			line := fmt.Sprintf("=== Line %d === ", lineNumber)
			n, err := writer.WriteString(line)
			if err != nil {
				log.Fatalf("Error writing: %v", err)
			}
			bytesWritten += int64(n)
		}

		// Write varied content patterns
		pattern := patterns[lineNumber%len(patterns)]
		n, err := writer.WriteString(pattern)
		if err != nil {
			log.Fatalf("Error writing: %v", err)
		}
		bytesWritten += int64(n)

		// Add newline every few patterns
		if lineNumber%3 == 0 {
			writer.WriteString("\n")
			bytesWritten++
		}

		lineNumber++

		// Prevent infinite loop
		if lineNumber > 1000000 {
			break
		}
	}

	// Ensure we end with a newline
	writer.WriteString("\n")

	fmt.Printf("Generated %s with %s of varied text content\n", filename, formatBytes(bytesWritten))
	fmt.Printf("Use: go run <source>.go -file=%s -cpuprofile=cpu.prof -memprofile=mem.prof\n", filename)
}

// parseCustomSize parses custom size specifications like "5MB", "500KB", "2GB"
func parseCustomSize(spec string) (int64, string) {
	if len(spec) < 3 {
		return 0, ""
	}

	// Extract number and unit
	var size int64
	var unit string

	if _, err := fmt.Sscanf(spec, "%d%s", &size, &unit); err != nil {
		return 0, ""
	}

	switch unit {
	case "B", "b":
		// size is already in bytes
	case "KB", "kb", "K", "k":
		size *= 1024
	case "MB", "mb", "M", "m":
		size *= 1024 * 1024
	case "GB", "gb", "G", "g":
		size *= 1024 * 1024 * 1024
	default:
		return 0, ""
	}

	// Generate filename based on size
	filename := fmt.Sprintf("/tmp/sample_custom_%s.txt", spec)

	return size, filename
}

func compareMethods(filename string) {
	// Profile unbuffered
	start := time.Now()
	unbufferedBytes := readUnbuffered(filename)
	unbufferedTime := time.Since(start)

	// Profile buffered
	start = time.Now()
	bufferedBytes := readBuffered(filename)
	bufferedTime := time.Since(start)

	// Show results
	fmt.Println("unbufferedBytes read:", unbufferedBytes)
	fmt.Println("bufferedBytes read:", bufferedBytes)
	fmt.Printf("Unbuffered: %v\n", unbufferedTime)
	fmt.Printf("Buffered:   %v\n", bufferedTime)
	fmt.Printf("Speedup:    %.2fx\n",
		float64(unbufferedTime)/float64(bufferedTime))
}
