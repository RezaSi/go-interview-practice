package main

import (
	"sync"
)

func bfs(graph map[int][]int, start int) []int {
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

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	if numWorkers < 1 {
		return map[int][]int{}
	}

	jobs := make(chan (int), numWorkers)
	go func() {
		defer close(jobs)
		for _, start := range queries {
			jobs <- start
		}
	}()

	var wg sync.WaitGroup
	wg.Add(len(queries))

	type result struct {
		start  int
		result []int
	}
	results := make(chan (*result), len(queries))
	for start := range jobs {
		go func() {
			defer wg.Done()
			res := &result{
				start:  start,
				result: bfs(graph, start),
			}
			results <- res
		}()
	}
	wg.Wait()
	close(results)

	r := map[int][]int{}
	for res := range results {
		r[res.start] = res.result
	}

	return r
}

func main() {
	// You can insert optional local tests here if desired.
}
