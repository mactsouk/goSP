package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

// Graph structure for adjacency map
type Graph struct {
	Nodes []string                      `json:"nodes"`
	Edges map[string]map[string]float64 `json:"edges"`
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func tspDP(g Graph) (float64, []string) {
	n := len(g.Nodes)
	index := make(map[string]int)
	for i, node := range g.Nodes {
		index[node] = i
	}

	size := 1 << n
	dp := make([][]float64, size)
	path := make([][]int, size)
	for i := range dp {
		dp[i] = make([]float64, n)
		path[i] = make([]int, n)
		for j := range dp[i] {
			dp[i][j] = math.Inf(1)
			path[i][j] = -1
		}
	}

	start := 0
	dp[1<<start][start] = 0

	for mask := 1; mask < size; mask++ {
		for u := 0; u < n; u++ {
			if mask&(1<<u) == 0 {
				continue
			}
			for v := 0; v < n; v++ {
				if mask&(1<<v) != 0 {
					continue
				}
				w, ok := g.Edges[g.Nodes[u]][g.Nodes[v]]
				if !ok {
					continue
				}
				newMask := mask | (1 << v)
				if dp[mask][u]+w < dp[newMask][v] {
					dp[newMask][v] = dp[mask][u] + w
					path[newMask][v] = u
				}
			}
		}
	}

	endMask := size - 1
	bestCost := math.Inf(1)
	last := -1
	for u := 0; u < n; u++ {
		w, ok := g.Edges[g.Nodes[u]][g.Nodes[start]]
		if !ok {
			continue
		}
		if dp[endMask][u]+w < bestCost {
			bestCost = dp[endMask][u] + w
			last = u
		}
	}

	if last == -1 {
		return math.Inf(1), nil
	}

	// reconstruct path
	route := make([]string, n)
	mask := endMask
	for i := n - 1; i >= 0; i-- {
		route[i] = g.Nodes[last]
		prev := path[mask][last]
		mask ^= 1 << last
		last = prev
	}

	return bestCost, route
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<json-file>")
		return
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var g Graph
	if err := json.Unmarshal(data, &g); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	cost, path := tspDP(g)
	if path != nil {
		fmt.Printf("Optimal path: %v\n", path)
		fmt.Printf("Total cost: %.2f\n", cost)
	} else {
		fmt.Println("No Hamiltonian cycle found")
	}
}
