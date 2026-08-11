package main

import (
	"fmt"
)

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.

type BFSResult struct {
	StartNode int
	Order     []int
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	if numWorkers <= 0 {
		return nil
	}
	jobs := make(chan int, len(queries))
	results := make(chan BFSResult, len(queries))
	// 创建任务
	for i := 0; i < numWorkers; i++ {
		go worker(graph, jobs, results)
	}
	// 分发任务
	for _, query := range queries {
		jobs <- query
	}
	close(jobs)
	// 收集结果
	orderedResults := make(map[int][]int, len(queries))
	for range queries {
		result := <-results
		orderedResults[result.StartNode] = result.Order
	}
	return orderedResults
}

func worker(graph map[int][]int, jobs <-chan int, results chan<- BFSResult) {
	for start := range jobs {
		order := query(graph, start)
		results <- BFSResult{StartNode: start, Order: order}
	}
}

func query(graph map[int][]int, start int) []int {
	queue := []int{start}
	visited := make(map[int]bool)
	var result []int
	visited[start] = true
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)
		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
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
		5: {2},
	}
	queries := []int{0, 1, 5}
	numWorkers := 2

	results := ConcurrentBFSQueries(graph, queries, numWorkers)
	fmt.Println(results)
}
