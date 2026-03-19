package main

import (
	"fmt"
	"strings"
)

type Node struct {
	prefix   string
	children []*Node
	isWord   bool
}

// RadixTree represents the root of the tree
type RadixTree struct {
	root *Node
}

// NewRadixTree creates a new Radix Tree
func NewRadixTree() *RadixTree {
	return new(RadixTree{
		root: new(Node{}),
	})
}

// Insert adds a word into the tree
func (tree *RadixTree) Insert(word string) {
	current := tree.root

	for {
		var child *Node
		var i int
		var j int
		for _, c := range current.children {
			// Find common prefix length
			j = commonPrefixLength(word, c.prefix)
			if j > 0 {
				i = j
				child = c
				break
			}
		}

		if child == nil {
			// No matching child, insert directly
			newNode := new(Node{prefix: word, isWord: true})
			current.children = append(current.children, newNode)
			return
		}

		// Found partial match or full match
		// i holds the length of the common prefix

		if i < len(child.prefix) {
			// Split the node
			// The existing child becomes a child of the new split node
			splitNode := new(Node{
				prefix:   child.prefix[i:],
				children: child.children,
				isWord:   child.isWord,
			})

			// Reset the child (which is now the parent of splitNode)
			child.prefix = child.prefix[:i]
			child.children = []*Node{splitNode}
			child.isWord = false
		}

		// Continue with the rest of the word
		word = word[i:]
		if word == "" {
			child.isWord = true
			return
		}
		current = child
	}
}

// Search checks if the word exists in the tree
func (tree *RadixTree) Search(word string) bool {
	current := tree.root

	for {
		var child *Node
		var i int
		for _, c := range current.children {
			if strings.HasPrefix(word, c.prefix) {
				child = c
				break
			}
			i = commonPrefixLength(word, c.prefix)
			if i > 0 {
				return false // Partial match means the word splits off here, so it's not found
			}
		}

		if child == nil {
			return false
		}

		if len(word) == len(child.prefix) {
			return child.isWord
		}

		word = word[len(child.prefix):]
		current = child
	}
}

// StartsWith checks if any word starts with the given prefix
func (tree *RadixTree) StartsWith(prefix string) bool {
	current := tree.root

	for {
		var child *Node
		for _, c := range current.children {
			if strings.HasPrefix(prefix, c.prefix) {
				child = c
				break
			}
			if strings.HasPrefix(c.prefix, prefix) {
				return true
			}
		}

		if child == nil {
			return false
		}

		if len(prefix) <= len(child.prefix) {
			return true
		}

		prefix = prefix[len(child.prefix):]
		current = child
	}
}

// commonPrefixLength returns the length of the common prefix
func commonPrefixLength(a, b string) int {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return i
}

// --- Test the RadixTree ---
func main() {
	tree := NewRadixTree()

	words := []string{"go", "gone", "good", "gopher", "goblin", "zebra"}
	for _, word := range words {
		tree.Insert(word)
	}

	tests := []string{"go", "gone", "god", "gob", "gopher", "zeb", "zebra"}
	for _, t := range tests {
		fmt.Printf("Search(%q): %v\n", t, tree.Search(t))
	}

	prefixes := []string{"go", "gop", "zeb", "x"}
	for _, p := range prefixes {
		fmt.Printf("StartsWith(%q): %v\n", p, tree.StartsWith(p))
	}
}
