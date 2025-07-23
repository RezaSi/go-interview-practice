package main

import "sync"

type Result struct {
	StartNode int
	BFSOrder  []int
}

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	var wg sync.WaitGroup
	queriesChan := make(chan int, len(queries))
	//resultChan := make(chan []int, len(queries))
	resultChan := make(chan Result, len(queries))

	result := make(map[int][]int)

	for i := range queries {
		queriesChan <- queries[i]
	}
	close(queriesChan)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for currentQuery := range queriesChan {
				worker(graph, currentQuery, resultChan)
			}
		}()

	}
	wg.Wait()
	close(resultChan)

	for k := range resultChan {
		result[k.StartNode] = k.BFSOrder
	}

	return result
}

func worker(graph map[int][]int, start int, ch chan Result) {
	visited := make(map[int]bool)
	query := make([]int, 0, 100)
	result := make([]int, 0, 100)
	query = append(query, start)
	visited[start] = true

	for len(query) > 0 {
		u := query[0]
		query = query[1:]
		result = append(result, u)
		for _, v := range graph[u] {
			if !visited[v] {
				visited[v] = true
				query = append(query, v)
			}
		}
	}
	ch <- Result{start, result}
}

func main() {
	// You can insert optional local tests here if desired.
}
