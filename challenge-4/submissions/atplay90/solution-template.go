package main

import (
	"fmt"
	"sync"
)

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
    
    if numWorkers <= 0 {
        return map[int][]int{} 
    }
	
	// Channels for distribution and collection
	jobs := make(chan int, len(queries))
	type result struct {
		start int
		order []int
	}
	resultsChan := make(chan result, len(queries))

	var wg sync.WaitGroup

	// 1. Spawn a fixed number of worker goroutines
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for start := range jobs {
				// Each worker processes a query sequentially
				order := bfs(graph, start)
				resultsChan <- result{start: start, order: order}
			}
		}()
	}

	// 2. Feed the jobs channel with queries
	for _, start := range queries {
		jobs <- start
	}
	close(jobs) // Closing signals workers to exit their range loops when done

	// 3. Wait for all workers to finish in a separate goroutine and close results
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 4. Collect results into the final map
	results := make(map[int][]int)
	for res := range resultsChan {
		results[res.start] = res.order
	}

	
	return results
}

func bfs(graph map[int][]int, start int) []int {
	// A simple reference BFS for checking correctness (sequential).
	queue := []int{start}
	visited := make(map[int]bool)
	visited[start] = true
	var order []int

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)

		for _, v := range graph[u] {
			if !visited[v] {
				visited[v] = true
				queue = append(queue, v)
			}
		}
	}

	return order
}

func main() {
	graph := map[int][]int{
        0: {1, 2},
        1: {2, 3},
        2: {3},
        3: {4},
        4: {},
    }
    queries := []int{0, 1, 2}
    numWorkers := 2

    results := ConcurrentBFSQueries(graph, queries, numWorkers)
    
    fmt.Println(results[0])
    fmt.Println(results[1])
    fmt.Println(results[2])
}
