package main

import (
	"fmt"
	"strings"
)

func main() {
	var builder strings.Builder
	// Preallocate to minimize memory allocations
	builder.Grow(128)

	// Write plain strings
	builder.WriteString("Go is ")
	builder.WriteString("fast!\n")

	// Append formatted data
	language := "Go"
	version := 1.25

	fmt.Fprintf(&builder, "Language: %s – ", language)
	fmt.Fprintf(&builder, "Version: %.2f\n", version)

	// More formatted examples
	cores := 12
	fmt.Fprintf(&builder, "Running on %d logical CPU cores\n", cores)

	// Output the final string
	final := builder.String()
	fmt.Print("Final output: ")
	fmt.Print(final)

	// Report internal statistics
	fmt.Printf("Length: %d bytes, Capacity: %d bytes\n", builder.Len(), builder.Cap())

	// Original string
	s := "Go is powerful!"

	// Count occurrences
	fmt.Println("Count of 's':", strings.Count(s, "s"))

	// Convert to uppercase and lowercase
	fmt.Println("Upper:", strings.ToUpper(s))
	fmt.Println("Lower:", strings.ToLower(s))

	// Replace words
	replaced := strings.ReplaceAll(s, "powerful", "fast")
	fmt.Println("REPLACED:", replaced)

	// Split into words
	words := strings.Fields(s)
	fmt.Println("Words:", words)

	// Join words with dashes
	dashed := strings.Join(words, "-")
	fmt.Println("Dashed:", dashed)

	// Index of a substring
	index := strings.Index(s, "powerful")
	fmt.Println("Index of 'powerful':", index)
}
