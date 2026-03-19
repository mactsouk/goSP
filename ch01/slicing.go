package main

import (
	"fmt"
)

func main() {
	// Define an array of 5 integers.
	arr := [5]int{10, 20, 30, 40, 50}
	fmt.Println("Original array:", arr)

	// Create two slices from the same array.
	slice1 := arr[1:4] // Elements at index 1, 2, 3 → [20, 30, 40]
	slice2 := arr[2:]  // Elements from index 2 to the end → [30, 40, 50]

	fmt.Println("Slice1:", slice1)
	fmt.Println("Slice2:", slice2)

	// Show that slices share the same underlying array.
	fmt.Println("\n--- Modifying slice1 ---")
	slice1[1] = 999
	fmt.Println("Slice1 after modification:", slice1)
	fmt.Println("Slice2 after modification:", slice2)
	fmt.Println("Array after modification:", arr)
	// Notice how slice2 and arr also reflect the change,
	// proving that they share the same underlying array.

	// Demonstrate slice length and capacity.
	fmt.Printf("\nSlice1 length: %d, capacity: %d\n", len(slice1), cap(slice1))

	// Now append within capacity.
	fmt.Println("\n--- Appending within capacity ---")
	slice1 = append(slice1, 888) // Still fits in the original array.
	fmt.Println("Slice1 after append:", slice1)
	fmt.Println("Array after append:", arr)
	// Since we haven't exceeded capacity, slice1 still points to arr.

	// Append beyond capacity — triggers allocation of a new array.
	fmt.Println("\n--- Appending beyond capacity ---")
	slice1 = append(slice1, 777) // Now we exceed the capacity.
	fmt.Println("Slice1 after exceeding capacity:", slice1)
	fmt.Println("Array after exceeding capacity:", arr)
	// At this point, slice1 no longer shares the underlying array with arr.

	// Verify independence after exceeding capacity.
	slice1[0] = 12345
	fmt.Println("Slice1 after independent modification:", slice1)
	fmt.Println("Array remains unchanged:", arr)
}
