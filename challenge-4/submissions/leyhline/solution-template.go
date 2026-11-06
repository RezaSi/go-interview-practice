package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func bfs(graph map[int][]int, query int) []int {
	var result []int
	result = append(result, query)
	currentDepth := graph[query]
	var nextDepth []int
	alreadyVisited := map[int]bool{query: true}
	for {
		if len(currentDepth) == 0 {
			return result
		}
		for _, vertice := range currentDepth {
			if alreadyVisited[vertice] {
				continue
			}
			result = append(result, vertice)
			alreadyVisited[vertice] = true
			nextDepth = append(nextDepth, graph[vertice]...)
		}
		currentDepth = nil
		currentDepth = append(currentDepth, nextDepth...)
		nextDepth = nil
	}
}

func runWorkers(graph map[int][]int, numWorkers int, queries []int) [][]int {
	results := make([][]int, len(queries))
	var next atomic.Int64
	var wg sync.WaitGroup
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				i := int(next.Add(1) - 1)
				if i >= len(queries) {
					return
				}
				results[i] = bfs(graph, queries[i])
			}
		}()
	}
	wg.Wait()
	return results
}

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	resultsMap := make(map[int][]int)
	if numWorkers <= 0 {
		return resultsMap
	}
	results := runWorkers(graph, numWorkers, queries)
	for i, query := range queries {
		resultsMap[query] = results[i]
	}
	return resultsMap
}

func main() {
	ownGraph := map[int][]int{
		1:  {2, 3},
		2:  {4, 5},
		3:  {6, 7},
		4:  {10, 11, 12},
		11: {31, 32},
		7:  {20, 21},
		31: {40, 2},
		32: {11, 1}}
	result := bfs(ownGraph, 1)
	fmt.Println("Result:", result)

	testGraph := map[int][]int{
		0: {1, 2},
		1: {2, 3},
		2: {3},
		3: {4},
		4: {},
		5: {2},
	}
	testResult := ConcurrentBFSQueries(testGraph, []int{0}, 1)
	fmt.Println("Result:", testResult)
}
