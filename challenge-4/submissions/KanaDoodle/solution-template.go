package main

import ("sync")
// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	if numWorkers <= 0 {
		return nil
	}
	
	type Result struct{
		start int
		order []int
	}

	bfs := func(start int) []int {
		visited := make(map[int]bool)
		order := []int{}
		queue := []int{start}
		visited[start] = true
		
		for len(queue) > 0 {
			node := queue[0]
			queue = queue[1:]
			order = append(order, node)
			for _, neighbor := range graph[node] {
				if !visited[neighbor] {
					visited[neighbor] = true
					queue = append(queue, neighbor)
				}
			}
		}
		return order
	}

	jobs := make(chan int)
	results := make(chan Result)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	worker := func ()  {
		defer wg.Done()
		for start := range jobs {
			order := bfs(start)
			results <- Result{start: start, order: order}
		}
	}

	for i := 0; i < numWorkers; i++ {
		go worker()
	}
	
	go func() {
		for _, q := range queries {
			jobs <- q
		}
		close(jobs)
	}()
	
	go func() {
		wg.Wait()
		close(results)
	}()

	finalResults := make(map[int][]int)
	for res := range results {
		finalResults[res.start] = res.order
	}

	return finalResults
}

func main(){
    
}