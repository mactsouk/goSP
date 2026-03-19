package main

import (
	"context"
	"fmt"
)

func main() {
	// Create a base context
	ctx := context.Background()

	// Derive a new context with key–value data
	ctx = context.WithValue(ctx, "userID", 42)
	ctx = context.WithValue(ctx, "role", "admin")
	processRequest(ctx)
}

func processRequest(ctx context.Context) {
	userID := ctx.Value("userID")
	role := ctx.Value("role")

	if userID != nil && role != nil {
		fmt.Printf("Processing request for user %v with role %v\n", userID, role)
	} else {
		fmt.Println("Missing user information in context")
	}
}
