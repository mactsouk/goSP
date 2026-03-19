package main

import (
	"fmt"
	"sync"
)

// Node represents a node in the binary search tree.
type Node struct {
	Value int
	Left  *Node
	Right *Node
}

// Tree represents the binary search tree itself.
// It includes a mutex to ensure thread safety.
type Tree struct {
	Root *Node
	mu   sync.RWMutex
}

// Insert inserts a value into the BST safely.
func (t *Tree) Insert(value int) {
	t.mu.Lock()         // Exclusive lock for writing
	defer t.mu.Unlock() // Unlock after insertion

	if t.Root == nil {
		t.Root = new(Node{Value: value})
	} else {
		t.Root.insert(value)
	}
}

// insert is the recursive helper for insertion.
func (n *Node) insert(value int) {
	if value <= n.Value {
		if n.Left == nil {
			n.Left = new(Node{Value: value})
		} else {
			n.Left.insert(value)
		}
	} else {
		if n.Right == nil {
			n.Right = new(Node{Value: value})
		} else {
			n.Right.insert(value)
		}
	}
}

// Search checks if a value exists in the tree safely.
func (t *Tree) Search(value int) bool {
	t.mu.RLock()         // Shared lock for reading
	defer t.mu.RUnlock() // Unlock after reading

	return t.Root != nil && t.Root.search(value)
}

// search is the recursive helper for searching.
func (n *Node) search(value int) bool {
	if n == nil {
		return false
	}
	if value == n.Value {
		return true
	} else if value < n.Value {
		return n.Left.search(value)
	} else {
		return n.Right.search(value)
	}
}

// InOrder prints the tree in sorted order safely.
func (t *Tree) InOrder() {
	t.mu.RLock()
	defer t.mu.RUnlock()

	fmt.Print("InOrder Traversal: ")
	inOrder(t.Root)
	fmt.Println()
}

func inOrder(n *Node) {
	if n == nil {
		return
	}
	inOrder(n.Left)
	fmt.Printf("%d ", n.Value)
	inOrder(n.Right)
}

// PreOrder prints the tree in pre-order safely.
func (t *Tree) PreOrder() {
	t.mu.RLock()
	defer t.mu.RUnlock()

	fmt.Print("PreOrder Traversal: ")
	preOrder(t.Root)
	fmt.Println()
}

func preOrder(n *Node) {
	if n == nil {
		return
	}
	fmt.Printf("%d ", n.Value)
	preOrder(n.Left)
	preOrder(n.Right)
}

// PostOrder prints the tree in post-order safely.
func (t *Tree) PostOrder() {
	t.mu.RLock()
	defer t.mu.RUnlock()

	fmt.Print("PostOrder Traversal: ")
	postOrder(t.Root)
	fmt.Println()
}

func postOrder(n *Node) {
	if n == nil {
		return
	}
	postOrder(n.Left)
	postOrder(n.Right)
	fmt.Printf("%d ", n.Value)
}

func main() {
	tree := new(Tree{})
	values := []int{8, 3, 10, 1, 6, 14, 4, 7, 13}
	var wg sync.WaitGroup

	for _, v := range values {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			tree.Insert(val)
			fmt.Printf("[Insert] %d\n", val)
		}(v)
	}
	wg.Wait()

	tree.InOrder()
	tree.PreOrder()
	tree.PostOrder()

	// Concurrent searches
	searchValues := []int{7, 15, 13, 2}
	wg = sync.WaitGroup{}

	for _, v := range searchValues {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			found := tree.Search(val)
			fmt.Printf("[Search] %d -> %v\n", val, found)
		}(v)
	}

	wg.Wait()
}
