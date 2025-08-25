package main

import (
	"fmt"
	"sync"
)

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

	if numWorkers <= 0 {
		return nil
	}

	bfs := func(query int) []int {
		visited := []int{}
		visitedSet := map[int]struct{}{}
		queue := []int{}
		queue = append(queue, query)
		for len(queue) > 0 {
			entry := queue[0]
			queue = queue[1:]
			if _, ok := visitedSet[entry]; ok {
				continue
			}
			visited = append(visited, entry)
			visitedSet[entry] = struct{}{}
			for _, adj := range graph[entry] {
				queue = append(queue, adj)
			}
		}

		return visited
	}

	type Task struct {
		ID  int
		Key int
	}

	type Result struct {
		ID      int
		Key     int
		Visited []int
	}

	chTask := make(chan Task)
	chResult := make(chan Result)
	results := map[int][]int{}

	wg := &sync.WaitGroup{}
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range chTask {
				chResult <- Result{
					ID:      task.ID,
					Key:     task.Key,
					Visited: bfs(task.Key),
				}
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for idx, query := range queries {
			chTask <- Task{ID: idx, Key: query}
		}
		close(chTask)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range len(queries) {
			result := <-chResult
			results[result.Key] = result.Visited
		}
		close(chResult)
	}()

	wg.Wait()

	return results
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
	fmt.Println(results)
}
