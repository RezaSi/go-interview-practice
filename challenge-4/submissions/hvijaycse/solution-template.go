package main

import "sync"

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.

type Result struct {
	start int
	order []int
}

func bfs(graph map[int][]int, start int) []int {

	order := []int{}
	visited := map[int]bool{}
	queue := []int{start}

	visited[start] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		order = append(order, node)

		for _, neighbor := range graph[node] {

			if visited[neighbor] {
				continue
			}
			visited[neighbor] = true
			queue = append(queue, neighbor)

		}
	}

	return order
}

func worker(wg *sync.WaitGroup, graph map[int][]int, jobs <-chan int, result chan<- Result) {

	defer wg.Done()

	for job := range jobs {
		order := bfs(graph, job)
		result <- Result{start: job, order: order}
	}

}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged.

	numWorkers = min(numWorkers, len(queries))
	wg := &sync.WaitGroup{}
	jobs := make(chan int)
	results := make(chan Result)

	output := make(map[int][]int)

	wg.Add(numWorkers)
	for _ = range numWorkers {
		go worker(wg, graph, jobs, results)
	}

	go func() {
		for _, query := range queries {
			jobs <- query
		}
		close(jobs)
	}()

	// Close results AFTER workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		output[result.start] = result.order
	}

	return output
}

func main() {
	// You can insert optional local tests here if desired.
}
