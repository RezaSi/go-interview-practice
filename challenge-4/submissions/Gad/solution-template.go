package main

import (
	"fmt"
	"maps"
	"slices"
)

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.

type Queue []int

func (q *Queue) Enqueue(node int) {
	*q = append(*q, node)
}

func (q *Queue) Dequeue() int {
	r := (*q)[0]
	*q = (*q)[1:]
	return r
}

func BFSQuery(graph map[int][]int, queryNode int) []int {

	var queue Queue
	var visited []int
	var results []int

	// edge case : empty graph - isolated node
	if len(graph) == 0 {
		return []int{queryNode}
	}

	// edge case : the node is not in the graph
	_, ok := graph[queryNode]
	if !ok {
		return []int{}
	}

	visited = append(visited, queryNode)
	queue.Enqueue(queryNode)

	for {
		if len(queue) != 0 {
			currentNode := queue.Dequeue()
			results = append(results, currentNode)
			neighbors := graph[currentNode]
			for _, n := range neighbors {
				if slices.Contains(visited, n) {
					continue
				}
				queue.Enqueue(n)
				visited = append(visited, n)
			}
		} else {
			break
		}
	}

	return results

}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {

	var (
		resultMap = make(map[int][]int)
		jobs      = make(chan int, len(queries))
		results   = make(chan map[int][]int, len(queries))
	)

	if numWorkers <= 0 {
		return nil
	}
	for range numWorkers {
		go func(chan int, chan map[int][]int) {
			for j := range jobs {
				results <- map[int][]int{
					j: BFSQuery(graph, j)}
			}
		}(jobs, results)
	}

	for _, j := range queries {
		jobs <- j
	}
	close(jobs)

	for range queries {
		maps.Copy(resultMap, <-results)
	}

	return resultMap
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

	fmt.Println(ConcurrentBFSQueries(graph, queries, numWorkers))
}
