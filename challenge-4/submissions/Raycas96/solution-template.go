package main

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.

type BFSResult struct {
	startNode int
	bfsOrder  []int
}

func standardBFS(graph map[int][]int, startNode int) []int {
    visited := make(map[int]bool)
    visited[startNode] = true
    queue := []int{startNode}
    order := []int{}

    for len(queue) > 0 {
        currentNode := queue[0]
        queue = queue[1:]
        
        order = append(order, currentNode)

        for _, neighbor := range graph[currentNode] {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
    return order 
}
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
    if numWorkers <= 0 {
        return map[int][]int{}
    }
    jobs := make(chan int, len(queries))
    results := make(chan BFSResult, len(queries))
	
	for _, query := range queries {
        jobs <- query
    }
    
    close(jobs)
    
    for i:=0; i< numWorkers; i++ {
        go func() {
            for startNode := range jobs {
                
                // "do stuff": run the standard BFS algorithm
                order := standardBFS(graph, startNode) 
                
                // Send the finished work back to the main function via the results channel
                results <- BFSResult{startNode: startNode, bfsOrder: order}
            }
        }()
}
	finalResults := make(map[int][]int)

    for i := 0; i < len(queries); i++ {
        res := <-results 
        
        finalResults[res.startNode] = res.bfsOrder
    }
    
    return finalResults
}

func main() {

}
