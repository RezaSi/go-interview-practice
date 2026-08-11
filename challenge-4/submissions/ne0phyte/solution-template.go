package main

import "sync"

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

	// channels for queries and results
	searchQueries := make(chan int, len(queries))
	searchResults := make(chan []int, len(queries))

	// spawn workers
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for range numWorkers {
		go func() {
			for query := range searchQueries {
				searchResults <- search(graph, query)
			}
			wg.Done()
		}()
	}

	// provide jobs (starting nodes of queries)
	for _, query := range queries {
		searchQueries <- query
	}
	close(searchQueries)

	// wait for workers and close jobsResults channel
	go func() {
		wg.Wait()
		close(searchResults)
	}()

	// collect results
	results := make(map[int][]int)
	for r := range searchResults {
		results[r[0]] = r
	}

	return results
}

func search(graph map[int][]int, start int) []int {
	visited := map[int]bool{start: true}
	queue := []int{start}
	result := []int{}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node)
		for _, adjacent := range graph[node] {
			if !visited[adjacent] {
				visited[adjacent] = true
				queue = append(queue, adjacent)
			}
		}
	}
	return result
}

func main() {
	// You can insert optional local tests here if desired.
}
