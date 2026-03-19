package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Note represents the structure of a note on the server
type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

const baseURL = "http://localhost:8080"

// --- Utility function: Pretty-print JSON response ---
func printJSON(body io.Reader) error {
	var data interface{}
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}
	pretty, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error formatting JSON: %w", err)
	}
	fmt.Println(string(pretty))
	return nil
}

// --- 1. List all notes ---
func listNotes() error {
	resp, err := http.Get(baseURL + "/notes")
	if err != nil {
		return fmt.Errorf("error fetching notes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("list failed with status: %s", resp.Status)
	}

	fmt.Println("Listing all notes:")
	return printJSON(resp.Body)
}

// --- 2. Get a specific note ---
func getNote(id int) error {
	resp, err := http.Get(fmt.Sprintf("%s/notes/%d", baseURL, id))
	if err != nil {
		return fmt.Errorf("error fetching note: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get note failed with status: %s", resp.Status)
	}

	fmt.Printf("Getting note with ID %d:\n", id)
	return printJSON(resp.Body)
}

// --- 3. Create a new note ---
func createNote(title, content string) error {
	note := Note{Title: title, Content: content}
	data, _ := json.Marshal(note)

	resp, err := http.Post(baseURL+"/notes", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creating note: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("create note failed with status: %s", resp.Status)
	}

	fmt.Println("Created new note successfully!")
	return printJSON(resp.Body)
}

// --- 4. Update an existing note ---
func updateNote(id int, title, content string) error {
	note := Note{Title: title, Content: content}
	data, _ := json.Marshal(note)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/notes/%d", baseURL, id), bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error updating note: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update note failed with status: %s", resp.Status)
	}

	fmt.Printf("Updated note with ID %d successfully!\n", id)
	return printJSON(resp.Body)
}

// --- 5. Delete a note ---
func deleteNote(id int) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/notes/%d", baseURL, id), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error deleting note: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete note failed with status: %s", resp.Status)
	}

	fmt.Printf("Deleted note with ID %d successfully.\n", id)
	return nil
}

// --- Main demo ---
func main() {
	fmt.Println("Go Note Client Demo")
	fmt.Println("---------------------")

	// 1. Create a couple of notes
	if err := createNote("First note", "This is the first note."); err != nil {
		fmt.Println(err)
	}
	if err := createNote("Second note", "This is the second note."); err != nil {
		fmt.Println(err)
	}

	// 2. List all notes
	if err := listNotes(); err != nil {
		fmt.Println(err)
	}

	// 3. Get a specific note
	if err := getNote(1); err != nil {
		fmt.Println(err)
	}

	// 4. Update note #1
	if err := updateNote(1, "Updated First Note", "This note has been updated."); err != nil {
		fmt.Println(err)
	}

	// 5. Delete note #2
	if err := deleteNote(2); err != nil {
		fmt.Println(err)
	}

	// 6. List all notes again
	if err := listNotes(); err != nil {
		fmt.Println(err)
	}
}

