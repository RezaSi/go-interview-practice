package main

import (
	"fmt"
	"sync"
)

// BFSResult represents the result of a single BFS query
type BFSResult struct {
	StartNode int
	Traversal []int
}

// BFS performs breadth-first search starting from the given node
func BFS(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	result := []int{}
	
	// Mark the starting node as visited
	visited[start] = true
	
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Explore all neighbors of current node (if it has any)
		if neighbors, exists := graph[current]; exists {
			for _, neighbor := range neighbors {
				if !visited[neighbor] {
					visited[neighbor] = true
					queue = append(queue, neighbor)
				}
			}
		}
	}

	return result
}


// ConcurrentBFSQueries processes BFS queries concurrently using worker goroutines
// Each worker processes queries from a shared job queue
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// Channel to distribute work to workers
	jobQueue := make(chan int, len(queries))
	
	// Channel to collect results from workers
	resultQueue := make(chan BFSResult, len(queries))
	
	// WaitGroup to ensure all workers complete before closing result channel
	var wg sync.WaitGroup
	
	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			// Each worker processes jobs until channel is closed
			for startNode := range jobQueue {
				bfsTraversal := BFS(graph, startNode)
				resultQueue <- BFSResult{
					StartNode: startNode,
					Traversal: bfsTraversal,
				}
			}
		}(i)
	}
	
	// Send all queries to the job queue
	for _, query := range queries {
		jobQueue <- query
	}
	close(jobQueue) // Signal no more jobs coming
	
	// Close result channel when all workers finish
	go func() {
		wg.Wait()
		close(resultQueue)
	}()
	
	// Collect all results
	finalResults := make(map[int][]int)
	for result := range resultQueue {
		finalResults[result.StartNode] = result.Traversal
	}
	
	return finalResults
}

// Alternative implementation using buffered channels for better performance
func ConcurrentBFSQueriesBuffered(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// Buffered channels for better throughput
	jobQueue := make(chan int, numWorkers*2)
	resultQueue := make(chan BFSResult, numWorkers*2)
	
	// Start workers
	for i := 0; i < numWorkers; i++ {
		go func() {
			for startNode := range jobQueue {
				bfsTraversal := BFS(graph, startNode)
				resultQueue <- BFSResult{
					StartNode: startNode,
					Traversal: bfsTraversal,
				}
			}
		}()
	}
	
	// Send jobs concurrently
	go func() {
		defer close(jobQueue)
		for _, query := range queries {
			jobQueue <- query
		}
	}()
	
	// Collect results
	finalResults := make(map[int][]int)
	for i := 0; i < len(queries); i++ {
		result := <-resultQueue
		finalResults[result.StartNode] = result.Traversal
	}
	
	return finalResults
}

func main() {
	// Example graph
	graph := map[int][]int{
		0: {1, 2},
		1: {2, 3},
		2: {3},
		3: {4},
		4: {},
	}
	
	queries := []int{0, 1, 2}
	numWorkers := 2

	fmt.Println("=== Standard Concurrent Implementation ===")
	results1 := ConcurrentBFSQueries(graph, queries, numWorkers)
	for _, query := range queries {
		fmt.Printf("BFS from %d: %v\n", query, results1[query])
	}
	
	fmt.Println("\n=== Buffered Channel Implementation ===")
	results2 := ConcurrentBFSQueriesBuffered(graph, queries, numWorkers)
	for _, query := range queries {
		fmt.Printf("BFS from %d: %v\n", query, results2[query])
	}
	
	// Larger example to demonstrate concurrency benefits
	fmt.Println("\n=== Large Graph Example ===")
	largeGraph := make(map[int][]int)
	for i := 0; i < 100; i++ {
		neighbors := []int{}
		if i < 99 {
			neighbors = append(neighbors, i+1)
		}
		if i > 0 {
			neighbors = append(neighbors, i-1)
		}
		largeGraph[i] = neighbors
	}
	
	largeQueries := []int{0, 25, 50, 75, 99}
	largeResults := ConcurrentBFSQueries(largeGraph, largeQueries, 4)
	
	for _, query := range largeQueries {
		result := largeResults[query]
		fmt.Printf("BFS from %d: [%d, %d, %d, ..., %d] (length: %d)\n", 
			query, result[0], result[1], result[2], result[len(result)-1], len(result))
	}
}