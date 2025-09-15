package main

import (
	"fmt"
	"sync"
)

// bfs runs a BFS traversal from a given start node.
func bfs(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	order := []int{}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if visited[node] {
			continue
		}
		visited[node] = true
		order = append(order, node)

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				queue = append(queue, neighbor)
			}
		}
	}
	return order
}

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	type task struct {
		start int
	}
	type result struct {
		start int
		order []int
	}

	tasks := make(chan task)
	results := make(chan result)

	var wg sync.WaitGroup

	// Worker pool
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range tasks {
				order := bfs(graph, t.start)
				results <- result{start: t.start, order: order}
			}
		}()
	}

	// Send tasks
	go func() {
		for _, q := range queries {
			tasks <- task{start: q}
		}
		close(tasks)
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	output := make(map[int][]int)
	for r := range results {
		output[r.start] = r.order
	}
	return output
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
	for q, order := range results {
		fmt.Printf("BFS from %d: %v\n", q, order)
	}
}
