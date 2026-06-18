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
	queue := make(chan int, len(queries))
	result := make(map[int][]int)

	for i := 0; i < len(queries); i++ {
		queue <- queries[i]
	}
	close(queue)

	var mutex sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for elem := range queue {
				mutex.Lock()
				result[elem] = bfs(graph, elem)
				mutex.Unlock()
			}
		}()
	}


	wg.Wait()

	return result
}

func bfs(graph map[int][]int, start int) []int {
	queue := []int{start}
	visited := make(map[int]bool)
	visited[start] = true
	var order []int

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)

		for _, v := range graph[u] {
			if !visited[v] {
				visited[v] = true
				queue = append(queue, v)
			}
		}
	}

	return order
}

func main() {
	// You can insert optional local tests here if desired.
}
