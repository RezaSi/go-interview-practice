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
    if numWorkers <= 0{
        // numWorkers = 1
        return map[int][]int{}
    }
    result := map[int][]int{}
    var resultMu sync.Mutex
    workerChan := make(chan struct{}, numWorkers)
    var wg sync.WaitGroup
    
    for _, query := range queries {
        wg.Add(1)
        workerChan <- struct{}{}
        
        go func(start int){
            defer wg.Done()
            defer func() { <-workerChan }()
            visited := make(map[int]bool)
            queue := []int{start}
            visited[start] = true
            bfs := []int{}
            
            for len(queue) > 0{
                node := queue[0]
                queue = queue[1:]
                bfs = append(bfs, node)
                
                for _, neghbour := range graph[node]{
                    if(visited[neghbour] != true){
                        queue = append(queue, neghbour)
                        visited[neghbour] = true
                    }
                }
            }
            resultMu.Lock()
            result[start] = bfs
            resultMu.Unlock()
        }(query)
    }
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged.
	wg.Wait()
    return result
}

func main() {
	// You can insert optional local tests here if desired.
}
