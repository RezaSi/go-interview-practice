package main

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged.

	type BFSResult struct {
		StartNode int
		Path      []int
	}

	jobs := make(chan int, len(queries))
	results := make(chan BFSResult, len(queries))

	for i := 0; i < numWorkers; i++ {
		go func() {
			for start := range jobs {
				path := bfs(graph, start)
				results <- BFSResult{
					StartNode: start,
					Path:      path,
				}
			}
		}()
	}

	for _, query := range queries {
		jobs <- query
	}
	close(jobs)

	if numWorkers == 0 {
		return map[int][]int{}
	}

	output := make(map[int][]int)
	for i := 0; i < len(queries); i++ {
		res := <-results
		output[res.StartNode] = res.Path
	}

	return output
}

func bfs(graph map[int][]int, start int) []int {
	queue := []int{start}
	var path []int
	vis := make(map[int]bool)

	vis[start] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		path = append(path, curr)

		for _, to := range graph[curr] {
			if vis[to] {
				continue
			}
			vis[to] = true
			queue = append(queue, to)
		}
	}

	return path
}
