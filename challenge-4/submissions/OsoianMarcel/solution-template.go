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
// Implementation notes:
//   - We spawn up to `numWorkers` goroutines that consume start-nodes from a
//     channel and run an independent BFS for each start node.
//   - Each BFS uses its own local `visited` map so there is no shared mutable
//     state during traversal. Results are sent back on a `results` channel and
//     collected by the caller goroutine which builds the final map.
//   - Worker goroutines are coordinated using a sync.WaitGroup; the results
//     channel is closed once all workers finish so the collector can range over it.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	if numWorkers <= 0 {
		return map[int][]int{}
	}

	tasks := make(chan int)
	type bfsResult struct {
		start int
		order []int
	}
	results := make(chan bfsResult)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// worker reads start nodes from tasks and performs BFS for each start
	worker := func() {
		defer wg.Done()
		for s := range tasks {
			// standard BFS using a slice as a queue
			visited := make(map[int]bool)
			order := make([]int, 0, 16)
			queue := make([]int, 0, 16)

			visited[s] = true
			queue = append(queue, s)

			for i := 0; i < len(queue); i++ {
				u := queue[i]
				order = append(order, u)
				for _, v := range graph[u] {
					if !visited[v] {
						visited[v] = true
						queue = append(queue, v)
					}
				}
			}

			// send the result back to collector
			results <- bfsResult{start: s, order: order}
		}
	}

	// start workers
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	// feed tasks
	go func() {
		for _, q := range queries {
			tasks <- q
		}
		close(tasks)
	}()

	// close results once all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// collect results into the output map
	out := make(map[int][]int, len(queries))
	for r := range results {
		out[r.start] = r.order
	}

	return out
}

func main() {
	// You can insert optional local tests here if desired.
}
