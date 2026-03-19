package main

import (
	"fmt"
	"math"
)

// 1. Define a regular interface
type Shape interface {
	Area() float64
	Perimeter() float64
}

// 2. Define concrete types
type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// 3. Generic function for Shape interface
func printShapeInfo(s Shape) {
	fmt.Printf("Area: %.2f ", s.Area())
	fmt.Printf("Perimeter: %.2f ", s.Perimeter())
	fmt.Println()
}

// 4. Function demonstrating the use of empty interfaces
func printAnything(v interface{}) {
	fmt.Printf("Value: %v, Type: %T\n", v, v)

	// Using a type assertion
	if str, ok := v.(string); ok {
		fmt.Println("This is a string with length:", len(str))
	}

	// Using a type switch for multiple types
	switch val := v.(type) {
	case int:
		fmt.Println("It's an integer, doubled:", val*2)
	case float64:
		fmt.Println("It's a float, squared:", val*val)
	case Shape:
		fmt.Println("It's a Shape:")
		printShapeInfo(val)
	default:
		fmt.Println("Unknown type!")
	}
}

func main() {
	rect := Rectangle{Width: 10, Height: 5}
	circle := Circle{Radius: 7}

	fmt.Println("Rectangle:")
	printShapeInfo(rect)

	fmt.Println("Circle:")
	printShapeInfo(circle)

	// Slice of Shape interfaces
	shapes := []Shape{
		Rectangle{3, 4},
		Circle{2.5},
		Rectangle{5, 8},
	}

	fmt.Println("Iterating over multiple shapes:")
	for i, s := range shapes {
		fmt.Printf("Shape %d:\n", i+1)
		printShapeInfo(s)
	}

	// 5. Using empty interfaces to hold mixed types
	fmt.Println("Working with empty interfaces:")
	items := []interface{}{
		"Go is awesome!",
		42,
		3.14,
		rect,
		circle,
		true,
	}

	for _, item := range items {
		printAnything(item)
	}
}
