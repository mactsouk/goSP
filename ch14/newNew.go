package main

import (
	"fmt"
)

type ServerConfig struct {
	Port    *int    `json:"port"`
	Enabled *bool   `json:"enabled"`
	Name    *string `json:"name"`
}

type User struct {
	ID   int
	Role string
}

func main() {
	// Before Go 1.26, you had to create variables to get pointers:
	// portVal := 8080
	// enabledVal := true
	// config := ServerConfig{Port: &portVal, Enabled: &enabledVal}

	// In Go 1.26, you can use new() with expressions directly:
	config := ServerConfig{
		Port:    new(8080),         // Returns *int initialized to 8080
		Enabled: new(true),         // Returns *bool initialized to true
		Name:    new("API Server"), // Returns *string initialized to "API Server"
	}

	// It also works with composite literals (structs, slices, arrays):
	// This allocates a User, initializes it, and returns *User.
	admin := new(User{
		ID:   1,
		Role: "admin",
	})

	fmt.Printf("Config: %s running on port %d (Enabled: %t)\n",
		*config.Name, *config.Port, *config.Enabled)
	fmt.Printf("User: %v\n", admin)
}
