package main

import "fmt"

// Node defines a node in a singly linked list
type Node struct {
	Value int
	Next  *Node
}

// LinkedList represents the linked list
type LinkedList struct {
	Head *Node
}

// Insert appends a new node to the end of the list
func (l *LinkedList) Insert(value int) {
	newNode := &Node{Value: value}
	if l.Head == nil {
		l.Head = newNode
		return
	}
	current := l.Head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode
}

// Delete removes the first node with the given value
func (l *LinkedList) Delete(value int) {
	if l.Head == nil {
		return
	}
	if l.Head.Value == value {
		l.Head = l.Head.Next
		return
	}
	prev := l.Head
	curr := l.Head.Next
	for curr != nil {
		if curr.Value == value {
			prev.Next = curr.Next
			return
		}
		prev = curr
		curr = curr.Next
	}
}

// Sort sorts the list in ascending order using bubble sort
func (l *LinkedList) Sort() {
	if l.Head == nil || l.Head.Next == nil {
		return
	}
	swapped := true
	for swapped {
		swapped = false
		curr := l.Head
		for curr.Next != nil {
			if curr.Value > curr.Next.Value {
				// Swap values
				curr.Value, curr.Next.Value = curr.Next.Value, curr.Value
				swapped = true
			}
			curr = curr.Next
		}
	}
}

// Print displays the values in the list
func (l *LinkedList) Print() {
	current := l.Head
	for current != nil {
		fmt.Print(current.Value, " -> ")
		current = current.Next
	}
	fmt.Println("nil")
}

func main() {
	list := &LinkedList{}

	// Insert values
	list.Insert(5)
	list.Insert(3)
	list.Insert(8)
	list.Insert(1)

	fmt.Println("Original list:")
	list.Print()

	// Delete a value
	list.Delete(3)
	fmt.Println("After deleting 3:")
	list.Print()

	// Sort the list
	list.Sort()
	fmt.Println("After sorting:")
	list.Print()
}
