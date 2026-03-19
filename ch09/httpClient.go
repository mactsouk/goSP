package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Helper function for making GET requests and printing responses
func getRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("[GET] %s\nStatus: %s\n", url, resp.Status)
	body, _ := io.ReadAll(resp.Body)
	fmt.Print("Response:", string(body))
}

// Helper function for making POST requests with JSON data
func postJSON(url string, data map[string]interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("[POST] %s\nStatus: %s\n", url, resp.Status)
	body, _ := io.ReadAll(resp.Body)
	fmt.Print("Response:", string(body))
}

// Helper function for making GET requests with custom headers
func getWithHeaders(url string, headers map[string]string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("[GET with Headers] %s\nStatus: %s\n", url, resp.Status)
	body, _ := io.ReadAll(resp.Body)
	fmt.Print("Response:", string(body))
}

func main() {
	server := "http://localhost:8080"
	if len(os.Args) > 1 {
		server = os.Args[1]
	}
	fmt.Printf("Connecting to server: %s", server)

	getRequest(server + "/")
	getRequest(server + "/time")
	getRequest(server + "/echo?msg=GoProgramming")
	postJSON(server+"/echo", map[string]interface{}{
		"user": "Mihalis",
		"lang": "Go",
	})
	getWithHeaders(server+"/headers", map[string]string{
		"X-Client": "GoHTTPClient",
		"X-Demo":   "true",
	})
	getRequest(server + "/health")
}
