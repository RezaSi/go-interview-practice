package main

import (
	"fmt"
	"sync"
)

func Bfs(graph map[int][]int, startingNode int) []int {
	visitedNodes := make(map[int]bool)
	visitedNodes[startingNode] = true
	queue := []int{startingNode}
	result := []int{}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)
		for _, neighbour := range graph[node] {
			if !visitedNodes[neighbour] {
				visitedNodes[neighbour] = true
				queue = append(queue, neighbour)
			}
		}
	}

	return result
}

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
	results := make(map[int][]int)
	var wg sync.WaitGroup
	var mu sync.Mutex
	channel := make(chan int)

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for startNode := range channel {
				bfsResult := Bfs(graph, startNode)
				mu.Lock()
				results[startNode] = bfsResult
				mu.Unlock()
			}
		}()
	}

	go func() {
		for _, query := range queries {
			channel <- query
		}
		close(channel)
	}()

	wg.Wait()
	return results
}

func main() {
	// You can insert optional local tests here if desired.
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
	fmt.Println(results)
}
