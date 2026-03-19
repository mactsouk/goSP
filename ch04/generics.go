package main

import (
	"fmt"
	"reflect"
)

// Define a custom constraint for numeric types
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Generic Sum function using our custom Numeric constraint
func Sum[T Numeric](nums []T) T {
	var total T
	for _, v := range nums {
		total += v
	}
	return total
}

// Non-generic Sum function using interface{}
func SumInterface(nums []interface{}) interface{} {
	if len(nums) == 0 {
		return nil
	}

	switch nums[0].(type) {
	case int:
		total := 0
		for _, v := range nums {
			total += v.(int)
		}
		return total
	case float64:
		total := 0.0
		for _, v := range nums {
			total += v.(float64)
		}
		return total
	case uint:
		var total uint = 0
		for _, v := range nums {
			total += v.(uint)
		}
		return total
	default:
		// fallback for unsupported types
		fmt.Printf("Unsupported: %v\n", reflect.TypeOf(nums[0]))
		return nil
	}
}

func main() {
	ints := []int{1, 2, 3, 4, 5}
	floats := []float64{1.5, 2.5, 3.5}
	uints := []uint{10, 20, 30}

	fmt.Println("Generics version:")
	fmt.Println("Sum of ints:", Sum(ints))     // works with []int
	fmt.Println("Sum of floats:", Sum(floats)) // works with []float64
	fmt.Println("Sum of uints:", Sum(uints))   // works with []uint

	fmt.Println("interface{} version:")
	fmt.Println("Sum of ints:", SumInterface([]interface{}{1, 2, 3, 4, 5}))
	fmt.Println("Sum of floats:", SumInterface([]interface{}{1.5, 2.5, 3.5}))
	fmt.Println("Sum of uints:", SumInterface([]interface{}{uint(10), uint(20), uint(30)}))
}
