package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var (
	notes      = make(map[int]Note)
	nextID     = 1
	notesMutex sync.Mutex
	dataFile   = "notes.json"
)

// Load notes from JSON file
func loadNotes() {
	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	defer file.Close()

	var loaded []Note
	if err := json.NewDecoder(file).Decode(&loaded); err != nil {
		panic(err)
	}

	notesMutex.Lock()
	defer notesMutex.Unlock()
	for _, n := range loaded {
		notes[n.ID] = n
		if n.ID >= nextID {
			nextID = n.ID + 1
		}
	}
}

// Save notes to JSON file
func saveNotes() {
	notesMutex.Lock()
	defer notesMutex.Unlock()

	allNotes := make([]Note, 0, len(notes))
	for _, n := range notes {
		allNotes = append(allNotes, n)
	}

	file, err := os.Create(dataFile)
	if err != nil {
		fmt.Println("Error saving notes:", err)
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(allNotes)
}

// JSON response helper
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Handlers
func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	var n Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	notesMutex.Lock()
	n.ID = nextID
	nextID++
	notes[n.ID] = n
	notesMutex.Unlock()

	saveNotes()
	writeJSON(w, http.StatusCreated, n)
}

func updateNoteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/notes/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	var n Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	notesMutex.Lock()
	if _, ok := notes[id]; !ok {
		notesMutex.Unlock()
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	n.ID = id
	notes[id] = n
	notesMutex.Unlock()

	saveNotes()
	writeJSON(w, http.StatusOK, n)
}

func deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/notes/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	notesMutex.Lock()
	if _, ok := notes[id]; !ok {
		notesMutex.Unlock()
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}
	delete(notes, id)
	notesMutex.Unlock()

	saveNotes()
	w.WriteHeader(http.StatusNoContent)
}

func listNotesHandler(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("title")

	notesMutex.Lock()
	defer notesMutex.Unlock()

	allNotes := make([]Note, 0, len(notes))
	for _, note := range notes {
		if filter == "" || strings.Contains(strings.ToLower(note.Title), strings.ToLower(filter)) {
			allNotes = append(allNotes, note)
		}
	}

	writeJSON(w, http.StatusOK, allNotes)
}

func getNoteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/notes/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	notesMutex.Lock()
	note, ok := notes[id]
	notesMutex.Unlock()
	if !ok {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, note)
}

// Middleware for logging
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		duration := time.Since(start)
		fmt.Printf("[%s] %s %s\nCompleted in %v\n", r.Method, r.URL.Path, r.RemoteAddr, duration)
	}
}

func main() {
	loadNotes()

	mux := http.NewServeMux()

	// /notes routes
	mux.HandleFunc("GET /notes", loggingMiddleware(listNotesHandler))
	mux.HandleFunc("POST /notes", loggingMiddleware(createNoteHandler))

	// /notes/{id} routes
	mux.HandleFunc("GET /notes/", loggingMiddleware(getNoteHandler))
	mux.HandleFunc("PUT /notes/", loggingMiddleware(updateNoteHandler))
	mux.HandleFunc("DELETE /notes/", loggingMiddleware(deleteNoteHandler))

	fmt.Println("Starting note-taking server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
