package main

import "sync"

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.

type bfsResult struct {
	start int
	order []int
}

func bfs(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	q := []int{start}
	order := []int{}

	for len(q) > 0 {
		node := q[0]
		q = q[1:]

		if visited[node] {
			continue
		}

		visited[node] = true
		order = append(order, node)

		for _, n := range graph[node] {
			if !visited[n] {
				q = append(q, n)
			}
		}
	}
	return order
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged.
	jobs := make(chan int)
	results := make(chan bfsResult)
	var wg sync.WaitGroup

	go func() {
		for _, q := range queries {
			jobs <- q
		}
		close(jobs)
	}()

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for start := range jobs {
				order := bfs(graph, start)
				results <- bfsResult{start, order}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	resMap := make(map[int][]int, len(queries))
	for r := range results {
		resMap[r.start] = r.order
	}

	return resMap
}

func main() {
	// You can insert optional local tests here if desired.
}
