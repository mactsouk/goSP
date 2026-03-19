package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// --- Protocol ---
type Message struct {
	Op    string `json:"op"`
	Host  string `json:"host,omitempty"`
	Path  string `json:"path,omitempty"`
	Query string `json:"query,omitempty"`
}

type Response struct {
	Status  string         `json:"status"`
	Results []string       `json:"results,omitempty"`
	Stats   map[string]int `json:"stats,omitempty"`
	Total   int            `json:"total_files,omitempty"`
	Message string         `json:"message,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "scan":
		runScan(os.Args[2:])
	case "watch":
		runWatch(os.Args[2:])
	case "search":
		runSearch(os.Args[2:])
	case "stats":
		runStats(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printHelp()
	}
}

func printHelp() {
	fmt.Println("Usage: client <command> [flags]")
	fmt.Println("Commands:")
	fmt.Println("  scan   - Index all files in a directory once")
	fmt.Println("  watch  - Monitor a directory for changes in real-time")
	fmt.Println("  search - Query the server for files")
	fmt.Println("  stats  - Show index statistics")
}

// --- Identity ---
func getMachineID() string {
	idFile := ".agent_id"

	content, err := os.ReadFile(idFile)
	if err == nil {
		return strings.TrimSpace(string(content))
	}

	fmt.Println("First run detected. Generating unique Agent ID...")
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}
	hostname = strings.ReplaceAll(hostname, ":", "_")

	fullID := fmt.Sprintf("%s-%s", hostname, hex.EncodeToString(bytes))
	os.WriteFile(idFile, []byte(fullID), 0644)
	return fullID
}

// --- COMMAND: SCAN ---
func runScan(args []string) {
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	dir := fs.String("dir", ".", "Directory to scan")
	server := fs.String("server", "localhost:8080", "Server address")
	fs.Parse(args)

	machineID := getMachineID()
	absDir, _ := filepath.Abs(*dir)

	fmt.Printf("Connecting to %s as '%s'...\n", *server, machineID)
	conn, err := net.Dial("tcp", *server)
	if err != nil {
		fmt.Printf("Connection error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	fmt.Printf("Scanning paths in %s...\n", absDir)
	count := 0
	start := time.Now()

	err = filepath.WalkDir(absDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}

		msg := Message{Op: "index", Host: machineID, Path: path}
		if err := encoder.Encode(msg); err != nil {
			return err
		}

		var resp Response
		if err := decoder.Decode(&resp); err != nil {
			return err
		}

		count++
		if count%1000 == 0 {
			fmt.Printf("\rIndexed %d files...", count)
		}
		return nil
	})

	fmt.Printf("\nDone. Sent %d paths in %v.\n", count, time.Since(start))
}

// --- COMMAND: WATCH ---
func runWatch(args []string) {
	fs := flag.NewFlagSet("watch", flag.ExitOnError)
	dir := fs.String("dir", ".", "Directory to watch")
	server := fs.String("server", "localhost:8080", "Server address")
	fs.Parse(args)

	machineID := getMachineID()
	absDir, _ := filepath.Abs(*dir)

	conn, err := net.Dial("tcp", *server)
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to server: %v", err))
	}
	defer conn.Close()
	encoder := json.NewEncoder(conn)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	fmt.Printf("Setting up watchers for %s...\n", absDir)
	dirCount := 0
	err = filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			dirCount++
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking path:", err)
	}
	fmt.Printf("Watching %d directories on machine '%s'.\n", dirCount, machineID)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			absPath, _ := filepath.Abs(event.Name)

			if event.Op&fsnotify.Create == fsnotify.Create {
				info, err := os.Stat(absPath)
				if err == nil && info.IsDir() {
					watcher.Add(absPath)
					fmt.Printf("[NEW DIR] %s\n", absPath)
				} else {
					fmt.Printf("[INDEX] %s\n", absPath)
					encoder.Encode(Message{Op: "index", Host: machineID, Path: absPath})
				}
			}

			if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
				fmt.Printf("[REMOVE] %s\n", absPath)
				encoder.Encode(Message{Op: "delete", Host: machineID, Path: absPath})
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Watcher error:", err)
		}
	}
}

// --- COMMAND: SEARCH ---
func runSearch(args []string) {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	query := fs.String("q", "", "Search term")
	server := fs.String("server", "localhost:8080", "Server address")
	fs.Parse(args)

	if *query == "" {
		fmt.Println("Please provide a query with -q")
		return
	}

	conn, err := net.Dial("tcp", *server)
	if err != nil {
		fmt.Printf("Server unreachable: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(Message{Op: "search", Query: *query}); err != nil {
		panic(err)
	}

	var resp Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		panic(err)
	}

	if len(resp.Results) == 0 {
		fmt.Printf("No results found for '%s'.\n", *query)
	} else {
		fmt.Printf("Found %d results:\n", len(resp.Results))
		for _, r := range resp.Results {
			parts := strings.SplitN(r, ":", 2)
			if len(parts) == 2 {
				fmt.Printf(" [%s] %s\n", parts[0], parts[1])
			} else {
				fmt.Printf(" - %s\n", r)
			}
		}
	}
}

// --- COMMAND: STATS ---
func runStats(args []string) {
	fs := flag.NewFlagSet("stats", flag.ExitOnError)
	server := fs.String("server", "localhost:8080", "Server address")
	fs.Parse(args)

	conn, err := net.Dial("tcp", *server)
	if err != nil {
		fmt.Printf("Server unreachable: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(Message{Op: "stats"}); err != nil {
		panic(err)
	}

	var resp Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		panic(err)
	}

	fmt.Println("========================================")
	fmt.Println("       Distributed Index Status         ")
	fmt.Println("========================================")
	fmt.Printf(" Total Files Indexed : %d\n", resp.Total)
	fmt.Println(" Active Agents       :")
	fmt.Println("----------------------------------------")

	// SORTING: Extract keys and sort them for consistent display
	var machines []string
	for k := range resp.Stats {
		machines = append(machines, k)
	}
	sort.Strings(machines)

	for _, host := range machines {
		count := resp.Stats[host]
		fmt.Printf(" %-25s : %d files\n", host, count)
	}
	fmt.Println("----------------------------------------")
}
