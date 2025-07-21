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
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// TODO: Implement concurrency-based BFS for multiple queries.
	type Task struct {
	    start int
	}
	
	type result struct {
	    start int
	    order []int
	}
	
	tasks := make(chan Task)
	results := make(chan result)
	
	var wg sync.WaitGroup
	
	for i := 0; i < numWorkers; i++ {
	    wg.Add(1)
	    go func() {
	        defer wg.Done()
	        for t := range tasks {
	            order := bfs(graph, t.start)
	            results <- result{start: t.start, order: order}
	        }
	    }()
	}
	
	go func() {
	    for _,start := range queries {
	        tasks <- Task{start: start}
	    }
	    close(tasks)
	}()
	
	go func() {
	    wg.Wait()
	    close(results)
	}()
	
	output := make(map[int][]int)
	for r := range results {
	    output[r.start] = r.order
	}
	// Return an empty map so the code compiles but fails tests if unchanged.
	return output
}

func bfs(grapfh map[int][]int, start int) []int {
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
        
        for _,neighbor := range grapfh[node] {
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
