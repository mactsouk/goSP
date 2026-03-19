package main

import "fmt"

func main() {
	x := 10
	fmt.Println("Before:", x)

	// Pass the address of x to the function
	updateValue(&x)
	fmt.Println("After:", x)
}

// Takes a pointer to an int as an argument
func updateValue(ptr *int) {
	*ptr = *ptr + 5 // Dereference the pointer and modify the value
}
