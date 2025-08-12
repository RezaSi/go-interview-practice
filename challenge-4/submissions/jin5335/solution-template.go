package main

import (
	"container/list"
	"fmt"
	"sync"
)

type BFSResult struct {
	StartNode int
	Order     []int
}

func worker(graph map[int][]int, jobCh chan int, resultCh chan BFSResult) {
	for startNode := range jobCh {
		q := list.New()
		visited := make(map[int]bool)
		bfsResult := BFSResult{StartNode: startNode, Order: make([]int, 0)}

		q.PushBack(startNode)
		visited[startNode] = true

		for q.Len() > 0 {

			var c *list.Element
			if c = q.Front(); c == nil {
				break
			}
			q.Remove(c)
			c_node := c.Value.(int)
			bfsResult.Order = append(bfsResult.Order, c_node)

			for _, neighbor := range graph[c_node] {
				if !visited[neighbor] {
					q.PushBack(neighbor)
					visited[neighbor] = true
				}
			}
		}
		resultCh <- bfsResult
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
	jobCh := make(chan int, len(queries))
	resultCh := make(chan BFSResult, len(queries))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(graph, jobCh, resultCh)
		}()
	}

	for _, query := range queries {
		jobCh <- query
	}
	close(jobCh)

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	result := make(map[int][]int, len(queries))
	for bfsResult := range resultCh {
		result[bfsResult.StartNode] = bfsResult.Order
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

	results := ConcurrentBFSQueries(graph, queries, numWorkers)
	for idx, r := range results {
		fmt.Printf("[#%d] %+v", idx, r)
	}
	/*
	   Possible output:
	   results[0] = [0 1 2 3 4]
	   results[1] = [1 2 3 4]
	   results[2] = [2 3 4]
	*/
	// You can insert optional local tests here if desired.
}
