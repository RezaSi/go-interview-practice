package main

import "sync"

// BFSResult encapsulates the result of a single BFS query.
// Using a struct is cleaner than sending a single-entry map over a channel.
type BFSResult struct {
	QueryNode int
	Order     []int
}

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// Use a channel for jobs (queries) to be sent to workers.
	queryChan := make(chan int, len(queries))
	// Use a channel for workers to send back their results.
	// The channel type is the result struct we defined.
	resultChan := make(chan BFSResult, len(queries))

	var wg sync.WaitGroup

	// Start the workers.
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Each worker ranges over the query channel until it's closed.
			for startNode := range queryChan {
				// Perform the BFS and send the result back.
				order := bfs(graph, startNode)
				resultChan <- BFSResult{QueryNode: startNode, Order: order}
			}
		}()
	}

	// Send all queries to the workers via the query channel.
	for _, query := range queries {
		queryChan <- query
	}
	// Close the query channel to signal to workers that no more jobs will be sent.
	close(queryChan)

	// Start a new goroutine that will close the resultChan once all workers are done.
	// This is the key to preventing deadlocks.

	wg.Wait()
	close(resultChan)


	// Collect all the results.
	// This loop will block until a result is available and will exit automatically
	// when resultChan is closed by the goroutine above.
	finalResults := make(map[int][]int, len(queries))
	for result := range resultChan {
		finalResults[result.QueryNode] = result.Order
	}

	return finalResults
}

// bfs performs a Breadth-First Search on a graph from a root node.
// This version is optimized to be more efficient.
func bfs(graph map[int][]int, root int)[]int{
    front := []int{root}
    seen := make(map[int]struct{})
    seen[root] = struct{}{}
    var resp []int
    for len(front) > 0{
        node := front[0]
        front = front[1:]
        resp = append(resp, node)
        for _, neighbour := range graph[node]{
            if _, ok := seen[neighbour]; !ok{
                seen[neighbour] = struct{}{}
                front = append(front, neighbour)
            }
        }
    }
    return resp
}

func main() {
	// You can insert optional local tests here if desired.
}