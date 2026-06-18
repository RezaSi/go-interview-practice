package main
import (
	"sync"
)
// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func bfs(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	visited[start] = true
	order := []int{}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	return order
}
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged.
	type job struct {
		query int
	}
	type result struct {
		query int
		order []int
	}

	jobs := make(chan job)
	results := make(chan result)

	var wg sync.WaitGroup

	// Worker function
	worker := func() {
		defer wg.Done()
		for j := range jobs {
			order := bfs(graph, j.query)
			results <- result{query: j.query, order: order}
		}
	}

	// Start workers
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// Send jobs
	go func() {
		for _, q := range queries {
			jobs <- job{query: q}
		}
		close(jobs)
	}()

	// Collect results in another goroutine
	done := make(chan struct{})
	resultMap := make(map[int][]int)
	var mu sync.Mutex

	go func() {
		for r := range results {
			mu.Lock()
			resultMap[r.query] = r.order
			mu.Unlock()
		}
		close(done)
	}()

	// Wait for all workers to finish
	wg.Wait()
	close(results)
	<-done

	return resultMap
}

func main() {
	// You can insert optional local tests here if desired.
}
