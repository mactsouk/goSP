package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println("=== Go Maps Comprehensive Example ===\n")

	// 1. Map Declaration and Initialization
	fmt.Println("1. Map Declaration and Initialization:")

	// Various ways to create maps
	var grades map[string]int      // Declare but not initialized (nil map)
	scores := make(map[string]int) // Create empty map with make()
	inventory := map[string]int{   // Map literal with initial values
		"apples": 50, "bananas": 30, "oranges": 25,
	}

	// Initialize the nil map
	grades = make(map[string]int)

	fmt.Printf("Empty scores map: %v\n", scores)
	fmt.Printf("Inventory map: %v\n", inventory)
	fmt.Println()

	// 2. Adding and Updating Elements
	fmt.Println("2. Adding and Updating Elements:")

	grades["Marietta"] = 95
	grades["Epifanios"] = 87
	grades["Athina"] = 92
	scores["Team A"] = 150
	scores["Team B"] = 200

	fmt.Printf("Grades after adding: %v\n", grades)
	fmt.Printf("Scores after adding: %v\n", scores)

	// Update existing value
	grades["Marietta"] = 98
	fmt.Printf("Updated Marietta's grade: %v\n", grades)
	fmt.Println()

	// 3. Reading Values and Zero Values
	fmt.Println("3. Reading Values and Zero Values:")

	mariettaGrade := grades["Marietta"]
	nonExistentGrade := grades["David"] // Returns zero value (0 for int)

	fmt.Printf("Marietta's grade: %d\n", mariettaGrade)
	fmt.Printf("David's grade (doesn't exist): %d\n", nonExistentGrade)
	fmt.Println()

	// 4. Checking if Key Exists (Comma OK Idiom)
	fmt.Println("4. Checking if Key Exists:")

	if grade, exists := grades["Marietta"]; exists {
		fmt.Printf("Marietta exists with grade: %d\n", grade)
	}

	if grade, exists := grades["David"]; exists {
		fmt.Printf("David exists with grade: %d\n", grade)
	} else {
		fmt.Printf("David doesn't exist in grades map\n")
	}
	fmt.Println()

	// 5. Deleting Elements
	fmt.Println("5. Deleting Elements:")

	fmt.Printf("Before deletion: %v\n", grades)
	delete(grades, "Epifanios")
	fmt.Printf("After deleting Epifanios: %v\n", grades)

	// Deleting non-existent key is safe (no-op)
	delete(grades, "NonExistent")
	fmt.Printf("After deleting non-existent key: %v\n", grades)
	fmt.Println()

	// 6. Iterating Over Maps
	fmt.Println("6. Iterating Over Maps:")

	fmt.Println("Iterating over key-value pairs:")
	for name, grade := range grades {
		fmt.Printf("  %s: %d\n", name, grade)
	}

	fmt.Println("Iterating over keys only:")
	for name := range grades {
		fmt.Printf("  Student: %s\n", name)
	}

	fmt.Println("Iterating over values only:")
	for _, grade := range grades {
		fmt.Printf("  Grade: %d\n", grade)
	}
	fmt.Println()

	// 7. Map Length and Capacity
	fmt.Println("7. Map Length:")

	fmt.Printf("Number of students: %d\n", len(grades))
	fmt.Printf("Number of inventory items: %d\n", len(inventory))
	fmt.Println()

	// 8. Maps with Different Key Types
	fmt.Println("8. Maps with Different Key Types:")

	// String keys
	cities := map[string]string{
		"NY": "New York",
		"CA": "California",
		"TX": "Texas",
	}

	// Integer keys
	fibonacci := map[int]int{
		1: 1, 2: 1, 3: 2, 4: 3, 5: 5, 6: 8,
	}

	// Boolean keys
	permissions := map[bool]string{true: "allowed", false: "denied"}

	fmt.Printf("Cities: %v\n", cities)
	fmt.Printf("Fibonacci: %v\n", fibonacci)
	fmt.Printf("Permissions: %v\n", permissions)
	fmt.Println()

	// 9. Nested Maps
	fmt.Println("9. Nested Maps:")

	studentGrades := map[string]map[string]int{
		"Alice": {"Math": 95, "Science": 92, "English": 88},
		"Bob":   {"Math": 87, "Science": 90, "English": 85},
	}

	fmt.Printf("Student grades: %v\n", studentGrades)
	fmt.Printf("Alice Math: %d\n", studentGrades["Alice"]["Math"])
	fmt.Println()

	// 10. Maps are Reference Types
	fmt.Println("10. Maps are Reference Types:")

	original := map[string]int{"a": 1, "b": 2}
	copy := original // This creates a reference, not a copy

	fmt.Printf("Original: %v\n", original)
	fmt.Printf("Copy: %v\n", copy)

	copy["c"] = 3 // Modifying copy affects original
	fmt.Printf("After modifying copy:\n")
	fmt.Printf("Original: %v\n", original)
	fmt.Printf("Copy: %v\n", copy)
	fmt.Println()

	// 11. Sorting Map Keys for Deterministic Iteration
	fmt.Println("11. Sorting Map Keys for Deterministic Iteration:")

	products := map[string]float64{
		"laptop": 999.99, "mouse": 29.99, "keyboard": 79.99, "monitor": 299.99,
	}

	// Get keys and sort them
	keys := make([]string, 0, len(products))
	for key := range products {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fmt.Println("Products sorted by name:")
	for _, key := range keys {
		fmt.Printf("  %s: $%.2f\n", key, products[key])
	}
	fmt.Println()

	// 12. Using Maps as Sets
	fmt.Println("12. Using Maps as Sets:")

	// Map with empty struct values to simulate a set
	uniqueWords := make(map[string]struct{})
	words := []string{"hello", "world", "hello", "go", "world", "programming"}

	for _, word := range words {
		uniqueWords[word] = struct{}{}
	}

	fmt.Printf("Original words: %v\n", words)
	fmt.Print("Unique words: ")
	for word := range uniqueWords {
		fmt.Printf("%s ", word)
	}
	fmt.Printf("\nNumber of unique words: %d\n", len(uniqueWords))
	fmt.Println()

	// 13. Safe Concurrent Access (Read-Only Example)
	fmt.Println("13. Map Operations Summary:")

	summary := map[string]interface{}{
		"total_students": len(grades),
		"average_grade":  calculateAverage(grades),
		"highest_grade":  findHighestGrade(grades),
		"map_size_bytes": "varies by implementation",
	}

	for key, value := range summary {
		fmt.Printf("  %s: %v\n", key, value)
	}
}

// Helper function to calculate average grade
func calculateAverage(grades map[string]int) float64 {
	if len(grades) == 0 {
		return 0
	}

	total := 0
	for _, grade := range grades {
		total += grade
	}

	return float64(total) / float64(len(grades))
}

// Helper function to find highest grade
func findHighestGrade(grades map[string]int) int {
	highest := 0
	for _, grade := range grades {
		if grade > highest {
			highest = grade
		}
	}
	return highest
}
