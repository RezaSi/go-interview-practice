package main

import (
	"container/list"
	"math"
	"sync"
)


func travel(graph map[int][]int, startNode int, result *sync.Map) {

	nodeCount := len(graph)

	visited := make([]bool, nodeCount, nodeCount)

	list := list.New()

	list.PushBack(startNode)
	visited[startNode] = true

	for list.Len() > 0 {

		node := list.Front().Value.(int)
		list.Remove(list.Front())

		l, ok := result.Load(startNode)
		if !ok {
			result.Store(startNode, []int{node})
		} else {
			result.Store(startNode, append(l.([]int), node))
		}

		for _, neigh := range graph[node] {

			if !visited[neigh] {
				visited[neigh] = true
				list.PushBack(neigh)

			}
		}
	}
}

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {

	if numWorkers == 0 || len(queries) == 0 {

		return map[int][]int{}
	}

	if len(graph) == 0 {
		finalResult := make(map[int][]int)
		for _, startNode := range queries {
			finalResult[startNode] = []int{startNode}
		}
		return finalResult
	}

	var result sync.Map

	var wg sync.WaitGroup

	for i := 1; i <= numWorkers; i++ {

		factor := int(math.Ceil(float64(len(queries)) / float64(numWorkers)))
		start := (i - 1) * factor
		end := min((i)*factor, len(queries))

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			for qIndex := start; qIndex < end; qIndex++ {

				travel(graph, queries[qIndex], &result)

			}

		}(start, end)
	}

	wg.Wait()

	finalResult := make(map[int][]int)

	for _, startNode := range queries {
		l, _ := result.Load(startNode)
		neigh, _ := l.([]int)
		finalResult[startNode] = neigh
	}

	return finalResult
}
func main() {
	// You can insert optional local tests here if desired.
}
