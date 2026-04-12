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
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged.
	if numWorkers == 0 {
		return map[int][]int{}
	}

	out := make(map[int][]int)
	var mu sync.Mutex
	var wg sync.WaitGroup
	tasks := make(chan int, len(queries))

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for index := range tasks {
				mu.Lock()
				out[index] = BFS(graph, index)
				mu.Unlock()
			}

		}()
	}

	for _, v := range queries {
		tasks <- v
	}

	close(tasks)

	wg.Wait()

	return out
}

func BFS(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	order := []int{}
	visited[start] = true

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		order = append(order, vertex)

		for _, neighbor := range graph[vertex] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return order

}

func main() {
	// You can insert optional local tests here if desired.
	graph := map[int][]int{
		0: {1, 2},
		1: {2, 3},
		2: {3},
		3: {4},
		4: {},
		5: {2},
	}
	queries := []int{0, 1, 5}

	numWorkers := 2

	results := ConcurrentBFSQueries(graph, queries, numWorkers)
	for k, v := range results {
		fmt.Println("X:", k, v)
	}

}
