package main

type BFSResult struct {
	StartNode int
	Order     []int
}

func bfs(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	visited[start] = true
	order := []int{}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		order = append(order, curr)

		for _, neighbour := range graph[curr] {
			if !visited[neighbour] {
				visited[neighbour] = true
				queue = append(queue, neighbour)
			}
		}
	}

	return order
}

func worker(graph map[int][]int, jobs <-chan int, results chan<- BFSResult) {
	for start := range jobs {
		order := bfs(graph, start)
		results <- BFSResult{StartNode: start, Order: order}
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
	jobs := make(chan int, len(queries))
	results := make(chan BFSResult, len(queries))
	resultMap := make(map[int][]int)

	if numWorkers <= 0 {
		close(jobs)
		close(results)
		return resultMap
	}

	for range numWorkers {
		go worker(graph, jobs, results)
	}

	for _, query := range queries {
		jobs <- query
	}
	close(jobs)

	for range len(queries) {
		result := <-results
		resultMap[result.StartNode] = result.Order
	}

	return resultMap
}

func main() {
	// You can insert optional local tests here if desired.
}
