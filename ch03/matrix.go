package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

func randomMatrix(rows, cols int) [][]int {
	mat := make([][]int, rows)
	for i := range mat {
		mat[i] = make([]int, cols)
		for j := range mat[i] {
			mat[i][j] = rand.Intn(3) - 1 // Random value: -1, 0, or 1
		}
	}
	return mat
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: ./<exe> <n> <m> <k>")
		os.Exit(1)
	}

	// Parse matrix dimensions
	n, _ := strconv.Atoi(os.Args[1])
	m, _ := strconv.Atoi(os.Args[2])
	k, _ := strconv.Atoi(os.Args[3])

	// Generate matrices
	A := randomMatrix(n, m)
	B := randomMatrix(m, k)

	// Initialize result matrix
	result := make([][]int, n)
	for i := range result {
		result[i] = make([]int, k)
	}

	// Perform matrix multiplication concurrently
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		for j := 0; j < k; j++ {
			wg.Add(1)
			go func(row, col int) {
				defer wg.Done()
				sum := 0
				for l := 0; l < m; l++ {
					sum += A[row][l] * B[l][col]
				}
				result[row][col] = sum // Safe: each goroutine writes to a unique cell
			}(i, j)
		}
	}
	wg.Wait()

	fmt.Println("Matrix A:")
	for _, row := range A {
		fmt.Println(row)
	}

	fmt.Println("Matrix B:")
	for _, row := range B {
		fmt.Println(row)
	}

	fmt.Println("Result Matrix (A x B):")
	for _, row := range result {
		fmt.Println(row)
	}
}
