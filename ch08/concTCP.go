package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// logWithTime prints log messages with a timestamp.
func logWithTime(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	fmt.Printf("[%s] %s\n", timestamp, message)
}

// handleConnection processes each client in its own goroutine.
func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	logWithTime("Client connected: %s", clientAddr)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		logWithTime("[%s] Received: %s", clientAddr, text)

		if strings.TrimSpace(text) == "quit" {
			logWithTime("[%s] Disconnected by client", clientAddr)
			return
		}

		// Echo message back
		response := fmt.Sprintf("Echo: %s\n", text)
		conn.Write([]byte(response))
	}

	if err := scanner.Err(); err != nil {
		logWithTime("[%s] Read error: %v", clientAddr, err)
	}

	logWithTime("Client %s connection closed", clientAddr)
}

func main() {
	// Define and parse the port flag
	port := flag.String("port", "8080", "TCP port number to listen on")
	flag.Parse()

	address := fmt.Sprintf(":%s", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listening on port %s: %v\n", *port, err)
		os.Exit(1)
	}
	defer listener.Close()

	logWithTime("Concurrent TCP server listening on port %s...", *port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logWithTime("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}
