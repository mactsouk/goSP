package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:", os.Args[0], "<number-of-nodes> <output-file>")
		return
	}

	var n int
	fmt.Sscanf(os.Args[1], "%d", &n)
	outputFile := os.Args[2]

	// Initialize nodes
	nodes := make([]string, n)
	for i := 0; i < n; i++ {
		nodes[i] = string('A' + i)
	}

	// Initialize edges map
	edges := make(map[string]map[string]float64)
	for _, u := range nodes {
		edges[u] = make(map[string]float64)
		for _, v := range nodes {
			if u == v {
				continue
			}
			edges[u][v] = math.Round(rand.Float64()*9 + 1) // weight 1..10
		}
	}

	graph := map[string]interface{}{
		"nodes": nodes,
		"edges": edges,
	}

	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		fmt.Println("Error generating JSON:", err)
		return
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Printf("Generated fully connected graph with %d nodes: %s\n", n, outputFile)
}
