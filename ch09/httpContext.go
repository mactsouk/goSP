package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Simulates a slow operation
func slowOperation(ctx context.Context) (string, error) {
	select {
	case <-time.After(3 * time.Second): // pretend work that takes 3 seconds
		return "Operation completed successfully!", nil
	case <-ctx.Done():
		return "", ctx.Err() // context canceled or timed out
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Create a new context with a 2-second timeout
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	result, err := slowOperation(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Request canceled: %v", err), http.StatusRequestTimeout)
		return
	}

	fmt.Fprintln(w, result)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
