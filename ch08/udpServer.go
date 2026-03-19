package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run udp_server.go <port>")
		return
	}

	port := os.Args[1]
	addr := net.UDPAddr{
		Port: mustParsePort(port),
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error starting UDP server:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("UDP server listening on port %s...\n", port)
	buffer := make([]byte, 1024)
	rand.Seed(time.Now().UnixNano())

	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			continue
		}

		fmt.Printf("Received %d bytes from %s: %s\n", n, clientAddr, string(buffer[:n]))

		// Generate random number and current time
		randomNum := rand.Intn(1000)
		currentTime := time.Now().Format(time.RFC1123)

		response := fmt.Sprintf("Random: %d | Time: %s\n", randomNum, currentTime)
		_, err = conn.WriteToUDP([]byte(response), clientAddr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}
	}
}

func mustParsePort(port string) int {
	p, err := net.LookupPort("udp", port)
	if err != nil {
		fmt.Println("Invalid port:", err)
		os.Exit(1)
	}
	return p
}
