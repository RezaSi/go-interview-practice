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

	type Result struct {
		start int
		order []int
	}

	jobChan := make(chan int)
	resultChan := make(chan Result)

	var wg sync.WaitGroup
	for range numWorkers {
		wg.Go(func() {
			for start := range jobChan {
				order := bfs(graph, start)
				resultChan <- Result{start, order}
			}
		})
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	go func() {
		for _, start := range queries {
			jobChan <- start
		}
		close(jobChan)
	}()

	results := make(map[int][]int)
	for result := range resultChan {
		results[result.start] = result.order
	}

	return results
}

// func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
// 	results := make(map[int][]int)
// 	if numWorkers <= 0 {
// 		return results
// 	}
// 	limitChan := make(chan struct{}, numWorkers)
// 	var mu sync.Mutex
// 	var wg sync.WaitGroup
// 	for _, start := range queries {
// 		limitChan <- struct{}{}
// 		wg.Go(func() {
// 			order := bfs(graph, start)
// 			mu.Lock()
// 			results[start] = order
// 			mu.Unlock()
// 			<-limitChan
// 		})
// 	}
// 	wg.Wait()
// 	return results
// }

func bfs(graph map[int][]int, start int) []int {
	var neighbors []int
	queue := []int{start}
	seen := map[int]bool{start: true}
	for n := len(queue); n > 0; n = len(queue) {
		for i := range n {
			node := queue[i]
			neighbors = append(neighbors, node)
			for _, nei := range graph[node] {
				if !seen[nei] {
					queue = append(queue, nei)
					seen[nei] = true
				}
			}
		}
		queue = queue[n:]
	}
	return neighbors
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
	for start, order := range results {
		fmt.Printf("%d -> %v\n", start, order)
	}
}
