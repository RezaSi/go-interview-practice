package main

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

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged.
    if numWorkers <= 0 {
        return map[int][]int{}
    }
    	
	jobs := make(chan int, len(queries))
	results := make(chan BFSResult, len(queries))

	for i := 0; i < numWorkers; i++ {
		go worker(graph, jobs, results)
	}

	for _, query := range queries {
		jobs <- query
	}
	close(jobs)

	resultMap := make(map[int][]int)

	for i := 0; i < len(queries); i++ {
		result := <-results
		resultMap[result.StartNode] = result.Order
	}

	return resultMap
}

func worker(graph map[int][]int, jobs <-chan int, results chan<- BFSResult) {
	for start := range jobs {
		order := bfs(graph, start)
		results <- BFSResult{StartNode: start, Order: order}
	}
}

func bfs(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	order := []int{}

	for len(queue) > 0 {

		node := queue[0]
		queue = queue[1:]

		if visited[node] {
			continue
		}

		visited[node] = true
		order = append(order, node)

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				queue = append(queue, neighbor)
			}
		}
	}

	return order
}

func main() {
	// You can insert optional local tests here if desired.
}
