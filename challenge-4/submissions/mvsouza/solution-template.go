package main

import (
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
	result := map[int][]int{}
	if numWorkers <= 0 || len(queries) == 0 {
		return result
	}

	var wg sync.WaitGroup
	tasks := make(chan int, len(queries))
	reducer := make(chan []int, len(queries))

	for i := range numWorkers {
		go func(id int) {
			for task := range tasks {
				reducer <- bfs(graph, task)
				wg.Done()
			}
		}(i)
	}

	for _, query := range queries {
		wg.Add(1)
		tasks <- query
	}
	close(tasks)
	go func() {
		wg.Wait()
		close(reducer)
	}()
	for reduce := range reducer {
		result[reduce[0]] = reduce
	}
	return result
}

func bfs(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	result := []int{}
	for len(queue) > 0 {
		var c int
		c, queue = queue[0], queue[1:]
		if visited[c] {
			continue
		}
		visited[c] = true
		queue = append(queue, graph[c]...)
		result = append(result, c)
	}
	return result
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

	ConcurrentBFSQueries(graph, queries, numWorkers)
}
