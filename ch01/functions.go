package main

import (
	"fmt"
	"strings"
)

// A simple function with parameters and a single return value.
func greet(name string) string {
	return "Hello, " + name + "!"
}

// A function returning multiple values.
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil
}

// A function with named return values.
func rectangleStats(length, width float64) (area, perimeter float64) {
	area = length * width
	perimeter = 2 * (length + width)
	return // Uses named returns
}

// A variadic function — accepts any number of string arguments.
func joinStrings(separator string, words ...string) string {
	return strings.Join(words, separator)
}

// A higher-order function: takes a function as an argument.
func applyTwice(fn func(int) int, value int) int {
	return fn(fn(value))
}

func main() {
	// 1. Basic function call.
	fmt.Println(greet("Mihalis"))

	// 2. Function returning multiple values.
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Division result:", result)
	}

	// 3. Named return values.
	area, perimeter := rectangleStats(5, 3)
	fmt.Printf("Rectangle → A: %.2f, P: %.2f\n", area, perimeter)

	// 4. Variadic function.
	joined := joinStrings("-", "Go", "is", "fast", "and", "fun")
	fmt.Println("Joined string:", joined)

	// 5. Higher-order function.
	double := func(x int) int { return x * 2 }
	fmt.Println("Applying double twice to 5:", applyTwice(double, 5))
}
