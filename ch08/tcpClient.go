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

// logWithTime prints timestamped log messages.
func logWithTime(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	fmt.Printf("[%s] %s\n", timestamp, message)
}

func main() {
	host := flag.String("host", "localhost", "Server")
	port := flag.String("port", "8080", "Server TCP port number")
	flag.Parse()

	address := fmt.Sprintf("%s:%s", *host, *port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to server at", address, err)
		os.Exit(1)
	}
	defer conn.Close()

	logWithTime("Connected to server at %s", address)
	fmt.Println("Type messages to send. Type 'quit' to exit.")

	// Channel to receive messages from the server
	go func() {
		serverScanner := bufio.NewScanner(conn)
		for serverScanner.Scan() {
			fmt.Printf("Server: %s\n", serverScanner.Text())
		}
		if err := serverScanner.Err(); err != nil {
			logWithTime("Error reading from server: %v", err)
		}
		os.Exit(0)
	}()

	// Read user input and send to the server
	inputScanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !inputScanner.Scan() {
			break
		}
		text := inputScanner.Text()

		_, err := fmt.Fprintln(conn, text)
		if err != nil {
			logWithTime("Error sending message: %v", err)
			break
		}

		if strings.TrimSpace(text) == "quit" {
			logWithTime("Disconnected from server.")
			break
		}
	}
}
