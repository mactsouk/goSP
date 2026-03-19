package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// JSON response helper
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Middleware for logging
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("[%s] %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		fmt.Printf("Completed in %v\n", time.Since(start))
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Go Web Service!")
}

func handleTime(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format(time.RFC3339)
	writeJSON(w, map[string]string{"time": now})
}

func handleEchoQuery(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	if msg == "" {
		msg = "No message provided"
	}
	writeJSON(w, map[string]string{"echo": msg})
}

func handleEchoPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]interface{}{"echo": body})
}

func handleHeaders(w http.ResponseWriter, r *http.Request) {
	headers := make(map[string][]string)
	for k, v := range r.Header {
		headers[k] = v
	}
	writeJSON(w, headers)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func main() {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", loggingMiddleware(handleRoot))
	mux.HandleFunc("GET /time", loggingMiddleware(handleTime))

	mux.HandleFunc("GET /echo", loggingMiddleware(handleEchoQuery))
	mux.HandleFunc("POST /echo", loggingMiddleware(handleEchoPost))

	mux.HandleFunc("GET /headers", loggingMiddleware(handleHeaders))
	mux.HandleFunc("GET /health", loggingMiddleware(handleHealth))

	fmt.Printf("Starting server on http://localhost:%s\n", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		panic(err)
	}
}
