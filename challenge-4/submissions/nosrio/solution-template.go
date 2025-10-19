package main

import "sync"

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

func bfs(graph map[int][]int, root int) []int {
	visited := make(map[int]bool)
	queue := []int{root}
	order := []int{}
	visited[root] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		order = append(order, current)

		for _, value := range graph[current] {
			if visited[value] {
				continue
			}
			visited[value] = true
			queue = append(queue, value)
		}
	}
	return order
}

func worker(wg *sync.WaitGroup, graph map[int][]int, jobs <-chan int, results chan<- BFSResult) {
	defer wg.Done()
	for start := range jobs {
		order := bfs(graph, start)
		results <- BFSResult{StartNode: start, Order: order}
	}
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged.
	resultMap := make(map[int][]int)
	var wg sync.WaitGroup

	if numWorkers == 0 {
		return resultMap
	}

	jobs := make(chan int, len(queries))
	results := make(chan BFSResult, len(queries))
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(&wg, graph, jobs, results)
	}

	for _, query := range queries {
		jobs <- query
	}
	close(jobs)
	wg.Wait()
	close(results)

	for i := 0; i < len(queries); i++ {
		result := <-results
		resultMap[result.StartNode] = result.Order
	}

	return resultMap
}

func main() {
	// You can insert optional local tests here if desired.
}
