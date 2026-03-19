package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println(os.Args[0], "<filename> <text to append>\n")
		os.Exit(1)
	}

	fn := os.Args[1]
	text := strings.Join(os.Args[2:], " ") + "\n"
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	_, err = f.WriteString(text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully appended content to %s\n", fn)
}
