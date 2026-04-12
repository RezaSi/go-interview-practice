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
	if graph == nil || numWorkers <= 0 || queries == nil{
	    return map[int][]int{}
	} 
	
	results := map[int][]int{}
	startChan := make(chan int, len(queries))
	resChan := make(chan []int, len(queries))
	
	for _, q := range queries{
	    startChan <-q
	}
	close(startChan)
	
	for i := 0; i < numWorkers; i++ {
	    go func(g map[int][]int, starts <-chan int, res chan<- []int) {
	        for start := range starts{
	            res<-BFS(g, start) 
	        }
	    }(graph, startChan, resChan)
	}
	
	for i:= 0; i < len(queries); i++{
	    result := <-resChan
	    results[result[0]] = result
	}
	return results
}

func BFS(graph map[int][]int, start int) []int{
    queue := []int{start}
    explored := map[int]bool{start: true}
    result := []int{start}
    
    for len(queue) > 0{
        curNode := queue[0]
        queue = queue[1:] //pop first node
        
        for _, node := range graph[curNode] {
            if !explored[node]{
                explored[node] = true
                result = append(result, node)
                queue = append(queue, node)
            }
        }
    }
    return result
}

func main() {
	// You can insert optional local tests here if desired.
}
