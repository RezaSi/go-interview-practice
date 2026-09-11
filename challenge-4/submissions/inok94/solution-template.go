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
type bfsResult struct {
	start  int
	result []int
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged.
	jobs := make(chan int)
	go func() {
		for _, v := range queries {
			jobs <- v
		}
		close(jobs)
	}()
	results := make(chan bfsResult)
	go workerPool(numWorkers, graph, jobs, results)
	res := make(map[int][]int)
	for r := range results {
		res[r.start] = r.result
	}
	return res
}

func bfs(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	result := []int{}

	visited[start] = true

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]

		result = append(result, vertex)

		for _, neigbor := range graph[vertex] {
			if !visited[neigbor] {
				visited[neigbor] = true
				queue = append(queue, neigbor)
			}
		}
	}
	return result
}

func workerPool(numsWorkers int, graph map[int][]int, jobs <-chan int, results chan<- bfsResult) {
	var wg sync.WaitGroup
	for i := 0; i < numsWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for start := range jobs {
				results <- bfsResult{start: start, result: bfs(graph, start)}
			}
		}()
	}
	wg.Wait()
	close(results)
}

func main() {
	// You can insert optional local tests here if desired.
}
