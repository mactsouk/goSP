package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// Job represents a unit of work.
type Job struct {
	ID int
}

// Result represents the outcome of processing a job.
type Result struct {
	JobID  int
	Output string
}

// worker processes jobs and sends results back.
func worker(id int, jobs <-chan Job, results chan<- Result) {
	for job := range jobs {
		// Simulate some work.
		fmt.Printf("Worker %d processing job %d\n", id, job.ID)
		time.Sleep(time.Millisecond * 500)

		// Send the result back.
		results <- Result{
			JobID:  job.ID,
			Output: fmt.Sprintf("Job %d – worker %d", job.ID, id),
		}
	}

	fmt.Printf("Worker %d exiting\n", id)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./<exe> <num_workers> <num_jobs>")
		os.Exit(1)
	}

	// Parse number of workers.
	numWorkers, err := strconv.Atoi(os.Args[1])
	if err != nil || numWorkers <= 0 {
		fmt.Println("Please provide a positive integer for the number of workers.")
		os.Exit(1)
	}

	// Parse number of jobs.
	numJobs, err := strconv.Atoi(os.Args[2])
	if err != nil || numJobs <= 0 {
		fmt.Println("Please provide a positive integer for the number of jobs.")
		os.Exit(1)
	}

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup

	// Start workers using WaitGroup.Go.
	fmt.Printf("%d workers for %d jobs...\n", numWorkers, numJobs)
	for w := 1; w <= numWorkers; w++ {
		workerID := w // safe in Go 1.25+, but explicit shadowing is still good practice
		wg.Go(func() {
			worker(workerID, jobs, results)
		})
	}

	// Send jobs to the workers in a separate goroutine.
	wg.Go(func() {
		for j := 1; j <= numJobs; j++ {
			jobs <- Job{ID: j}
		}
		close(jobs)
	})

	// Wait for all workers and close results channel when done.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and print results.
	for result := range results {
		fmt.Println("Result:", result.Output)
	}

	fmt.Println("All jobs processed!")
}
