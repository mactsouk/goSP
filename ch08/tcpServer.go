package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	// Listen on TCP port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Non-concurrent TCP server listening on port 8080...")

	for {
		// Wait for a client to connect
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Printf("Client: %s\n", conn.RemoteAddr().String())
		// Handle the client directly (no goroutine)
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("Received: %s\n", text)
		if strings.TrimSpace(text) == "quit" {
			fmt.Println("Close connection request.")
			return
		}
		// Echo message back to client
		response := fmt.Sprintf("Echo: %s\n", text)
		conn.Write([]byte(response))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from connection:", err)
	}
}
