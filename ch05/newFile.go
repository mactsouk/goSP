package main

import (
	"fmt"
	"os"
)

func main() {
	fn := "/tmp/unique.txt"

	// os.O_CREATE → create file if not exists
	// os.O_EXCL   → fail if file already exists
	// os.O_WRONLY → open for writing
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			fmt.Printf("File %s already exists.\n", fn)
			return
		}
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	content := "This file was created only because it did not exist.\n"
	_, err = f.WriteString(content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("File %s created successfully.\n", fn)
}
