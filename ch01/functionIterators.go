package main

import (
	"fmt"
)

// fibonacci returns a function iterator that lazily
// generates Fibonacci numbers.
func fibonacci(limit int) func(yield func(int) bool) {
	return func(yield func(int) bool) {
		a, b := 0, 1
		for i := 0; i < limit; i++ {
			if !yield(a) {
				return // Stop early if the caller stops consuming values
			}
			a, b = b, a+b
		}
	}
}

func main() {
	fmt.Println("First 10 Fibonacci numbers:")

	// Iterate lazily over the Fibonacci sequence
	for num := range fibonacci(10) {
		fmt.Print(num, " ")
	}
	fmt.Println()
}
