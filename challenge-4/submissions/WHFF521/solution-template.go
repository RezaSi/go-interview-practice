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
type result struct {
    query int
    order []int
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged.
	numQueries := len(queries)
	jobs :=make(chan int, numQueries)
	results := make(chan result, numQueries)
	var wg sync.WaitGroup
	
	for i:=0; i < numWorkers; i++ {
	    wg.Add(1)
	    go func() {
	        defer wg.Done()
	        for startNode := range jobs {
	            bfsOrder := bfs(graph, startNode)
	            results <- result{query:startNode,order:bfsOrder}
	        }
	    }()
	}
	
	for _, q := range queries {
	    jobs <-q
	}
	close(jobs)
	
	go func(){
	    wg.Wait()
	    close(results)
	}()
	
	finalResults := make(map[int][]int)
	for res :=range results {
	    finalResults[res.query] = res.order
	}
	
	return finalResults
}

func bfs(graph map[int][]int, start int) []int {
    order := []int{}
    visited := make(map[int]bool)
    queue := []int{start}
    visited[start] = true
    
    for len(queue) > 0 {
        curr := queue[0]
        queue = queue[1:]
        order = append(order , curr)
        
        for _,neighbor := range graph[curr]{
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue,neighbor)
            }
        }
    }
    return order
}

func main() {
	// You can insert optional local tests here if desired.
}
