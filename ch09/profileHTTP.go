package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// JSON response helper
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Middleware using fmt.Printf()
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
	port := flag.String("port", "8080", "Port to run the HTTP server on")
	flag.Parse()

	var middleware func(http.HandlerFunc) http.HandlerFunc
	middleware = loggingMiddleware

	// Main application routes
	http.HandleFunc("/", middleware(handleRoot))
	http.HandleFunc("/time", middleware(handleTime))
	http.HandleFunc("/echo", middleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handleEchoQuery(w, r)
		} else if r.Method == http.MethodPost {
			handleEchoPost(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	http.HandleFunc("/headers", middleware(handleHeaders))
	http.HandleFunc("/health", middleware(handleHealth))

	// Profiling server
	go func() {
		log.Println("Profiler starting at http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Println("pprof server error:", err)
		}
	}()

	log.Printf("Starting server on http://localhost:%s", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		log.Fatal(err)
	}
}
