package main

import (
	"container/heap"
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

// Item for priority queue
type Item struct {
	node     string
	priority float64
	index    int
}

// PriorityQueue implements heap.Interface
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// Dijkstra algorithm
func dijkstra(g Graph, source string) (map[string]float64, map[string]string) {
	dist := make(map[string]float64)
	prev := make(map[string]string)
	for _, node := range g.Nodes {
		dist[node] = math.Inf(1)
		prev[node] = ""
	}
	dist[source] = 0

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Item{node: source, priority: 0})

	for pq.Len() > 0 {
		u := heap.Pop(pq).(*Item)
		for v, w := range g.Edges[u.node] {
			alt := dist[u.node] + w
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = u.node
				heap.Push(pq, &Item{node: v, priority: alt})
			}
		}
	}

	return dist, prev
}

// Reconstruct path from source to target
func reconstructPath(prev map[string]string, target string) []string {
	path := []string{}
	for u := target; u != ""; u = prev[u] {
		path = append([]string{u}, path...)
	}
	return path
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:", os.Args[0], "<json-file> <source-node>")
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

	source := os.Args[2]
	dist, prev := dijkstra(g, source)

	fmt.Printf("Shortest distances from %s:\n", source)
	for _, node := range g.Nodes {
		fmt.Printf("%s: %.2f\n", node, dist[node])
	}

	fmt.Println("\nPaths from source:")
	for _, node := range g.Nodes {
		if node == source {
			continue
		}
		path := reconstructPath(prev, node)
		fmt.Printf("%s -> %s\n", source, node)
		fmt.Printf("Path: %v\n", path)
	}
}
