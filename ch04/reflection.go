package main

import (
	"fmt"
	"reflect"
)

// A sample struct to demonstrate reflection
type User struct {
	Name string
	Age  int
}

// A function to be called dynamically using reflection
func Greet(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

func main() {
	// --- 1. Inspecting basic types ---
	fmt.Println("=== Inspecting Types and Values ===")
	var x float64 = 3.14

	v := reflect.ValueOf(x)
	t := reflect.TypeOf(x)
	fmt.Println("Type:", t)        // float64
	fmt.Println("Value:", v)       // 3.14
	fmt.Println("Kind:", v.Kind()) // float64

	// --- 2. Modifying struct fields dynamically ---
	fmt.Println("=== Modifying Struct Fields ===")
	u := User{Name: "Alice", Age: 30}
	fmt.Println("Before:", u)

	// Get a reflect.Value of the struct pointer
	rv := reflect.ValueOf(&u).Elem()

	// Modify the "Name" field
	nameField := rv.FieldByName("Name")
	if nameField.IsValid() && nameField.CanSet() {
		nameField.SetString("Bob")
	}

	// Modify the "Age" field
	ageField := rv.FieldByName("Age")
	if ageField.IsValid() && ageField.CanSet() {
		ageField.SetInt(42)
	}
	fmt.Println("After:", u)

	// --- 3. Calling a function dynamically ---
	fmt.Println("=== Calling Functions Dynamically ===")
	greetFunc := reflect.ValueOf(Greet)
	args := []reflect.Value{reflect.ValueOf("Charlie")}
	greetFunc.Call(args)
}
