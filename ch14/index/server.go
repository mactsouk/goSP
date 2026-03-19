package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

// --- Protocol ---
type Message struct {
	Op    string `json:"op"`             // "index", "delete", "search", "stats"
	Host  string `json:"host,omitempty"` // "MachineA-123"
	Path  string `json:"path,omitempty"` // "/var/www/index.html"
	Query string `json:"query,omitempty"`
}

type Response struct {
	Status  string         `json:"status"`
	Results []string       `json:"results,omitempty"`
	Stats   map[string]int `json:"stats,omitempty"`
	Total   int            `json:"total_files,omitempty"`
	Message string         `json:"message,omitempty"`
}

// --- Data Structures ---
// Node represents a character in the Trie
type Node struct {
	Children map[string]*Node
	Files    []string // Stores "Host:Path" strings
}

// Journal handles the Append-Only Log
type Journal struct {
	file *os.File
	mu   sync.Mutex
}

// Server holds the state of the entire application
type Server struct {
	Root     *Node
	Registry map[string]bool // Keeps track of unique "Host:Path" entries
	Stats    map[string]int  // Aggregated counts: "MachineA" -> 500 files
	Total    int             // Total indexed files
	Journal  *Journal
	mu       sync.RWMutex
}

// --- Journal Logic (Persistence) ---
func OpenJournal(filename string) (*Journal, error) {
	// Open file in Append mode, Create if not exists, Write Only
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Journal{file: f}, nil
}

func (j *Journal) WriteEvent(msg Message) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	// Add newline delimiter
	if _, err := j.file.Write(append(data, '\n')); err != nil {
		return err
	}
	return j.file.Sync()
}

func (j *Journal) Close() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.file.Close()
}

// --- Server Logic ---
func NewServer(journalPath string) *Server {
	j, err := OpenJournal(journalPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to open journal: %v", err))
	}

	s := &Server{
		Root:     &Node{Children: make(map[string]*Node)},
		Registry: make(map[string]bool),
		Stats:    make(map[string]int),
		Journal:  j,
	}

	// Replay history immediately upon creation
	s.replayJournal(journalPath)
	return s
}

func (s *Server) replayJournal(path string) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		} // No journal yet, that's fine
		panic(err)
	}
	defer file.Close()

	fmt.Println("Replaying journal to restore index...")
	scanner := bufio.NewScanner(file)
	count := 0
	start := time.Now()

	for scanner.Scan() {
		var msg Message
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			continue
		}
		if msg.Op == "index" {
			s.IndexInternal(msg.Host, msg.Path)
		} else if msg.Op == "delete" {
			s.DeleteInternal(msg.Host, msg.Path)
		}
		count++
	}
	fmt.Printf("Restored %d events in %v. Total files: %d\n", count, time.Since(start), s.Total)
}

func tokenize(path string) []string {
	cleanPath := strings.ReplaceAll(path, "\\", "/")
	return strings.FieldsFunc(cleanPath, func(r rune) bool {
		return r == '/' || r == '.' || r == '_' || r == '-'
	})
}

// IndexInternal updates the Trie and Stats (Thread-Unsafe)
func (s *Server) IndexInternal(host, path string) {
	fullID := fmt.Sprintf("%s:%s", host, path)

	if s.Registry[fullID] {
		return
	}
	s.Registry[fullID] = true
	s.Stats[host]++
	s.Total++

	tokens := tokenize(path)
	for _, token := range tokens {
		if len(token) < 2 {
			continue
		}
		token = strings.ToLower(token)

		node := s.Root
		for _, char := range token {
			charStr := string(char)
			if node.Children == nil {
				node.Children = make(map[string]*Node)
			}
			if _, exists := node.Children[charStr]; !exists {
				node.Children[charStr] = &Node{Children: make(map[string]*Node)}
			}
			node = node.Children[charStr]
		}
		exists := false
		for _, f := range node.Files {
			if f == fullID {
				exists = true
				break
			}
		}
		if !exists {
			node.Files = append(node.Files, fullID)
		}
	}
}

// DeleteInternal removes from Trie and Stats (Thread-Unsafe)
func (s *Server) DeleteInternal(host, path string) {
	fullID := fmt.Sprintf("%s:%s", host, path)

	if !s.Registry[fullID] {
		return
	}
	delete(s.Registry, fullID)
	s.Stats[host]--
	if s.Stats[host] <= 0 {
		delete(s.Stats, host)
	}
	s.Total--

	tokens := tokenize(path)
	for _, token := range tokens {
		if len(token) < 2 {
			continue
		}
		token = strings.ToLower(token)

		node := s.Root
		for _, char := range token {
			if next, ok := node.Children[string(char)]; ok {
				node = next
			} else {
				node = nil
				break
			}
		}

		if node != nil {
			newFiles := make([]string, 0)
			for _, f := range node.Files {
				if f != fullID {
					newFiles = append(newFiles, f)
				}
			}
			node.Files = newFiles
		}
	}
}

// Public Methods (Thread-Safe)
func (s *Server) HandleIndex(host, path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Journal.WriteEvent(Message{Op: "index", Host: host, Path: path})
	s.IndexInternal(host, path)
}

func (s *Server) HandleDelete(host, path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Journal.WriteEvent(Message{Op: "delete", Host: host, Path: path})
	s.DeleteInternal(host, path)
}

func (s *Server) HandleSearch(query string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	node := s.Root
	for _, char := range strings.ToLower(query) {
		if next, ok := node.Children[string(char)]; ok {
			node = next
		} else {
			return nil
		}
	}

	var results []string
	seen := make(map[string]bool)
	var collect func(*Node)
	collect = func(n *Node) {
		for _, file := range n.Files {
			if !seen[file] {
				results = append(results, file)
				seen[file] = true
			}
		}
		for _, child := range n.Children {
			collect(child)
		}
	}
	collect(node)

	// SORTING: Ensure results are alphabetical
	sort.Strings(results)

	return results
}

func (s *Server) HandleStats() (int, map[string]int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	statsCopy := make(map[string]int)
	for k, v := range s.Stats {
		statsCopy[k] = v
	}
	return s.Total, statsCopy
}

// --- Connection Handler ---
func handleConnection(conn net.Conn, server *Server) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var msg Message
		if err := decoder.Decode(&msg); err == io.EOF {
			return
		} else if err != nil {
			return
		}

		switch msg.Op {
		case "index":
			server.HandleIndex(msg.Host, msg.Path)
			encoder.Encode(Response{Status: "ok"})
		case "delete":
			server.HandleDelete(msg.Host, msg.Path)
			encoder.Encode(Response{Status: "deleted"})
		case "search":
			results := server.HandleSearch(msg.Query)
			encoder.Encode(Response{Status: "ok", Results: results})
		case "stats":
			total, breakdown := server.HandleStats()
			encoder.Encode(Response{Status: "ok", Total: total, Stats: breakdown})
		}
	}
}

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	server := NewServer("server.wal")

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Printf("\nReceived signal: %s. Shutting down...\n", sig)
		listener.Close()
		server.Journal.Close()
		fmt.Println("Journal closed. Goodbye.")
		os.Exit(0)
	}()

	fmt.Printf("Running on port %s (PID: %d)\n", *port, os.Getpid())
	fmt.Println("Press Ctrl+C to shutdown.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-sigs:
				return
			default:
				fmt.Printf("Connection error: %v\n", err)
			}
			continue
		}
		go handleConnection(conn, server)
	}
}
