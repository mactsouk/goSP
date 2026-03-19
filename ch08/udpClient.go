package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run udp_client.go <host> <port>")
		return
	}
	host := os.Args[1]
	port := os.Args[2]

	// Combine host and port into "host:port" format
	address := net.JoinHostPort(host, port)

	// Resolve the UDP address
	serverAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	// Create a UDP connection to the resolved address
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	message := "request"
	fmt.Printf("Sending message to %s: %s\n", address, message)

	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	buffer := make([]byte, 1024)
	// Prevent blocking forever
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Error receiving response:", err)
		return
	}

	fmt.Printf("Server response: %s\n", string(buffer[:n]))
}
