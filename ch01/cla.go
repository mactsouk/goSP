package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("=== Raw Command-Line Arguments ===")

	// Use filepath.Base() to get only the program's basename.
	programName := filepath.Base(os.Args[0])
	fmt.Printf("Program name: %s\n", programName)

	// Display all arguments except the program name.
	if len(os.Args) > 1 {
		fmt.Println("Arguments:")
		for i, arg := range os.Args[1:] {
			fmt.Printf("  Arg %d: %s\n", i+1, arg)
		}
	} else {
		fmt.Println("No additional arguments provided.")
	}

	fmt.Println("=== Using the flag Package ===")

	// Define flags for structured argument parsing.
	port := flag.Int("port", 8080, "Port to listen on")
	host := flag.String("host", "localhost", "Hostname or IP")
	debug := flag.Bool("debug", false, "Enable debug mode")

	// Parse the flags from os.Args.
	flag.Parse()

	fmt.Printf("Server will start at %s:%d\n", *host, *port)
	fmt.Printf("Debug mode: %v\n", *debug)

	// Show remaining non-flag arguments, if any.
	remaining := flag.Args()
	if len(remaining) > 0 {
		fmt.Println("Non-flag arguments:", remaining)
	} else {
		fmt.Println("No non-flag arguments provided.")
	}
}
