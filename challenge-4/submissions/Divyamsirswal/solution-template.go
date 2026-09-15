package main

import (
	"fmt"
)

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {

	if len(queries) == 0 {
		return make(map[int][]int)
	}

	if numWorkers <= 0 {
		return make(map[int][]int)
	}

	type Job struct{ start int }
	type Result struct {
		start int
		nodes []int
	}

	jobs := make(chan Job, len(queries))
	results := make(chan Result, len(queries))

	for w := 0; w < numWorkers; w++ {

		go func() {

			for job := range jobs {
				x := []int{}
				queue := []int{job.start}
				vis := make(map[int]bool)
				vis[job.start] = true

				for len(queue) != 0 {
					node := queue[0]
					queue = queue[1:]
					x = append(x, node)

					for _, neigh := range graph[node] {
						if vis[neigh] == false {
							vis[neigh] = true
							queue = append(queue, neigh)
						}
					}
				}

				results <- Result{start: job.start, nodes: x}
			}

		}()

	}

	for _, q := range queries {
		jobs <- Job{start: q}
	}
	close(jobs)

	res := make(map[int][]int)
	for i := 0; i < len(queries); i++ {
		r := <-results
		res[r.start] = r.nodes
	}

	return res

}

func main() {
	graph := map[int][]int{
		0: {1, 2},
		1: {2, 3},
		2: {3},
		3: {4},
		4: {},
	}
	queries := []int{0, 1, 2}
	numWorkers := 2

	results := ConcurrentBFSQueries(graph, queries, numWorkers)
	/*
	   Possible output:
	   results[0] = [0 1 2 3 4]
	   results[1] = [1 2 3 4]
	   results[2] = [2 3 4]
	*/

	fmt.Println(results)
}
