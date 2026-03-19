package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

// 1. Basic Struct Definition
type Person struct {
	Name string
	Age  int
	City string
}

// 2. Struct with Different Field Types
type Product struct {
	ID        int
	Name      string
	Price     float64
	InStock   bool
	Tags      []string
	CreatedAt time.Time
}

// 3. Struct with Embedded Fields (Composition)
type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
}

type Employee struct {
	Person     // Embedded struct (anonymous field)
	Address    // Embedded struct
	ID         int
	Department string
	Salary     float64
	IsManager  bool
}

// 4. Struct with Methods (Receiver Functions)
type Rectangle struct {
	Width  float64
	Height float64
}

// Value receiver - doesn't modify the struct
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Value receiver
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Pointer receiver - can modify the struct
func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

// Pointer receiver
func (r *Rectangle) SetDimensions(width, height float64) {
	r.Width = width
	r.Height = height
}

// 5. Struct with Private and Public Fields
type BankAccount struct {
	AccountNumber string  // Public (exported)
	balance       float64 // Private (unexported)
	ownerName     string  // Private
}

// Constructor function (common Go pattern)
func NewBankAccount(accountNumber, ownerName string, initialBalance float64) *BankAccount {
	return &BankAccount{
		AccountNumber: accountNumber,
		balance:       initialBalance,
		ownerName:     ownerName,
	}
}

// Getter method for private field
func (b *BankAccount) Balance() float64 {
	return b.balance
}

// Setter method for private field with validation
func (b *BankAccount) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive")
	}
	b.balance += amount
	return nil
}

func (b *BankAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdrawal amount must be positive")
	}
	if amount > b.balance {
		return fmt.Errorf("insufficient funds")
	}
	b.balance -= amount
	return nil
}

// 6. Interface Implementation
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14159 * c.Radius
}

// 7. Struct with Tags (for JSON, database mapping, etc.)
type User struct {
	ID       int    `json:"id" db:"user_id" validate:"required"`
	Username string `json:"username" db:"username" validate:"required,min=3"`
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"-" db:"password_hash"`          // "-" excludes from JSON
	FullName string `json:"full_name,omitempty" db:"name"` // omitempty excludes if empty
}

// 8. Anonymous Structs
type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
	Server struct {
		Port int    `json:"port"`
		Host string `json:"host"`
	} `json:"server"`
}

// 9. Struct Comparison and Copying
type Point struct {
	X, Y int
}

func main() {
	fmt.Println("=== Go Structures Comprehensive Example ===\n")

	// 1. Basic Struct Usage
	fmt.Println("1. Basic Struct Declaration and Initialization:")

	// Different ways to create structs
	var p1 Person                         // Zero value initialization
	p2 := Person{}                        // Empty initialization
	p3 := Person{"Alice", 30, "New York"} // Positional initialization
	p4 := Person{                         // Named field initialization
		Name: "Bob",
		Age:  25,
		City: "Boston",
	}

	fmt.Printf("p1 (zero value): %+v\n", p1)
	fmt.Printf("p2 (empty): %+v\n", p2)
	fmt.Printf("p3 (positional): %+v\n", p3)
	fmt.Printf("p4 (named): %+v\n", p4)
	fmt.Println()

	// 2. Accessing and Modifying Fields
	fmt.Println("2. Accessing and Modifying Fields:")

	p1.Name = "Charlie"
	p1.Age = 35
	p1.City = "Chicago"

	fmt.Printf("Modified p1: %+v\n", p1)
	fmt.Printf("p1.Name: %s, p1.Age: %d\n", p1.Name, p1.Age)
	fmt.Println()

	// 3. Struct Pointers
	fmt.Println("3. Struct Pointers:")

	p5 := &Person{"David", 40, "Denver"} // Pointer to struct
	fmt.Printf("p5 pointer: %p\n", p5)
	fmt.Printf("p5 value: %+v\n", *p5)

	// Go automatically dereferences struct pointers
	fmt.Printf("p5.Name (auto-dereference): %s\n", p5.Name)

	// Explicit dereferencing – not recommended but available
	fmt.Printf("(*p5).Name (explicit): %s\n", (*p5).Name)

	// Modifying through pointer
	p5.Age = 41
	fmt.Printf("Modified through pointer: %+v\n", *p5)
	fmt.Println()

	// 4. Embedded Structs (Composition)
	fmt.Println("4. Embedded Structs (Composition):")

	emp := Employee{
		Person: Person{
			Name: "Emily",
			Age:  28,
			City: "Seattle",
		},
		Address: Address{
			Street:  "123 Main St",
			City:    "Seattle",
			State:   "WA",
			ZipCode: "98101",
		},
		ID:         1001,
		Department: "Engineering",
		Salary:     75000,
		IsManager:  false,
	}

	fmt.Printf("Employee: %+v\n", emp)

	// Accessing embedded fields directly
	fmt.Printf("Employee name: %s\n", emp.Name)     // From embedded Person
	fmt.Printf("Employee street: %s\n", emp.Street) // From embedded Address
	fmt.Printf("Employee department: %s\n", emp.Department)

	// Accessing embedded fields explicitly
	fmt.Printf("Person info: %+v\n", emp.Person)
	fmt.Printf("Address info: %+v\n", emp.Address)
	fmt.Println()

	// 5. Methods on Structs
	fmt.Println("5. Methods on Structs:")

	rect := Rectangle{Width: 10, Height: 5}
	fmt.Printf("Rectangle: %+v\n", rect)
	fmt.Printf("Area: %.2f\n", rect.Area())
	fmt.Printf("Perimeter: %.2f\n", rect.Perimeter())

	// Method with pointer receiver
	fmt.Printf("Before scaling: %+v\n", rect)
	rect.Scale(2)
	fmt.Printf("After scaling by 2: %+v\n", rect)

	rect.SetDimensions(15, 8)
	fmt.Printf("After setting new dimensions: %+v\n", rect)
	fmt.Println()

	// 6. Private Fields and Constructor Pattern
	fmt.Println("6. Private Fields and Constructor Pattern:")

	account := NewBankAccount("ACC-001", "John Doe", 1000.0)
	fmt.Printf("New account: %+v\n", account)
	fmt.Printf("Initial balance: $%.2f\n", account.Balance())

	// Deposit money
	err := account.Deposit(500)
	if err != nil {
		fmt.Printf("Deposit error: %v\n", err)
	} else {
		fmt.Printf("After deposit: $%.2f\n", account.Balance())
	}

	// Withdraw money
	err = account.Withdraw(300)
	if err != nil {
		fmt.Printf("Withdrawal error: %v\n", err)
	} else {
		fmt.Printf("After withdrawal: $%.2f\n", account.Balance())
	}

	// Try invalid operations
	err = account.Withdraw(2000)
	if err != nil {
		fmt.Printf("Invalid withdrawal: %v\n", err)
	}
	fmt.Println()

	// 7. Interface Implementation
	fmt.Println("7. Interface Implementation:")

	shapes := []Shape{
		Rectangle{Width: 10, Height: 5},
		Circle{Radius: 3},
	}

	for i, shape := range shapes {
		fmt.Printf("Shape %d:\n", i+1)
		fmt.Printf("  Type: %T\n", shape)
		fmt.Printf("  Area: %.2f\n", shape.Area())
		fmt.Printf("  Perimeter: %.2f\n", shape.Perimeter())
	}
	fmt.Println()

	// 8. Struct Tags and JSON
	fmt.Println("8. Struct Tags and JSON:")

	user := User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "secret123",
		FullName: "John Doe",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("JSON marshal error: %v\n", err)
	} else {
		fmt.Printf("JSON representation: %s\n", string(jsonData))
	}

	// Reflect on struct tags
	userType := reflect.TypeOf(user)
	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)
		jsonTag := field.Tag.Get("json")
		dbTag := field.Tag.Get("db")
		fmt.Printf("Field %s: json='%s', db='%s'\n", field.Name, jsonTag, dbTag)
	}
	fmt.Println()

	// 9. Anonymous Structs
	fmt.Println("9. Anonymous Structs:")

	// Inline anonymous struct
	settings := struct {
		Theme    string
		Language string
		Debug    bool
	}{
		Theme:    "dark",
		Language: "en",
		Debug:    true,
	}

	fmt.Printf("Settings: %+v\n", settings)

	// Anonymous struct in slice
	people := []struct {
		Name string
		Role string
	}{
		{"Alice", "Developer"},
		{"Bob", "Designer"},
		{"Charlie", "Manager"},
	}

	fmt.Println("People:")
	for _, person := range people {
		fmt.Printf("  %s - %s\n", person.Name, person.Role)
	}
	fmt.Println()

	// 10. Nested Anonymous Structs
	fmt.Println("10. Nested Anonymous Structs:")

	config := Config{}
	config.Database.Host = "localhost"
	config.Database.Port = 5432
	config.Database.Username = "admin"
	config.Database.Password = "secret"
	config.Server.Host = "0.0.0.0"
	config.Server.Port = 8080

	fmt.Printf("Config: %+v\n", config)
	fmt.Println()

	// 11. Struct Comparison and Copying
	fmt.Println("11. Struct Comparison and Copying:")

	point1 := Point{X: 1, Y: 2}
	point2 := Point{X: 1, Y: 2}
	point3 := Point{X: 3, Y: 4}

	fmt.Printf("point1: %+v\n", point1)
	fmt.Printf("point2: %+v\n", point2)
	fmt.Printf("point3: %+v\n", point3)
	fmt.Printf("point1 == point2: %t\n", point1 == point2)
	fmt.Printf("point1 == point3: %t\n", point1 == point3)

	// Copying structs (value semantics)
	point4 := point1 // Creates a copy
	point4.X = 10
	fmt.Printf("Original point1: %+v\n", point1)
	fmt.Printf("Modified copy point4: %+v\n", point4)
	fmt.Println()

	// 12. Struct with Slice and Map Fields
	fmt.Println("12. Struct with Complex Fields:")

	product := Product{
		ID:        101,
		Name:      "Laptop",
		Price:     999.99,
		InStock:   true,
		Tags:      []string{"electronics", "computer", "portable"},
		CreatedAt: time.Now(),
	}

	fmt.Printf("Product: %+v\n", product)
	fmt.Printf("Tags: %s\n", strings.Join(product.Tags, ", "))
	fmt.Printf("Created: %s\n", product.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()

	// 13. Empty Struct Usage
	fmt.Println("13. Empty Struct Usage:")

	// Empty struct as signal
	type Signal struct{}

	// Using empty struct as set values (memory efficient)
	visited := make(map[string]struct{})
	items := []string{"apple", "banana", "apple", "cherry", "banana"}

	for _, item := range items {
		visited[item] = struct{}{}
	}

	fmt.Printf("Original items: %v\n", items)
	fmt.Print("Unique items: ")
	for item := range visited {
		fmt.Printf("%s ", item)
	}
	fmt.Printf("\nEmpty struct size: %d bytes\n", unsafe.Sizeof(struct{}{}))
	fmt.Println()

	// 14. Method Sets and Value vs Pointer Receivers
	fmt.Println("14. Method Sets Summary:")

	rect1 := Rectangle{Width: 5, Height: 3}
	rect2 := &Rectangle{Width: 7, Height: 4}

	// Both value and pointer can call value receiver methods
	fmt.Printf("rect1.Area(): %.2f\n", rect1.Area())
	fmt.Printf("rect2.Area(): %.2f\n", rect2.Area())

	// Both value and pointer can call pointer receiver methods
	// (Go automatically takes address of value or dereferences pointer)
	rect1.Scale(1.5)
	rect2.Scale(2.0)

	fmt.Printf("After scaling - rect1: %+v\n", rect1)
	fmt.Printf("After scaling - rect2: %+v\n", *rect2)
}
