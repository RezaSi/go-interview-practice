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
    var wg sync.WaitGroup
    var m sync.Mutex
    ch := make(chan struct{}, numWorkers)
    result := make(map[int][]int, len(queries))
    if numWorkers == 0 {
        return result
    }
    maxVertex := 0
    for _, query := range queries {
        maxVertex = max(maxVertex, query)
    }
    wg.Add(len(queries))
    for _, startNode := range queries {
        go func(start int) {
            ch<-struct{}{}
            var d []int
            used := make(map[int]struct{})
            var queue []int
            queue = append(queue, start)
            for len(queue) > 0 {
                v := queue[0]
                queue = queue[1:]
                if _, found := used[v]; found {
                    continue
                }
                used[v] = struct{}{}
                d = append(d, v)
                for _, to := range graph[v] {
                    if _, found := used[to]; !found {
                        queue = append(queue, to)
                    }
                }
            }
            m.Lock()
            result[start] = d
            m.Unlock()
            wg.Done()
            <-ch
        }(startNode)   
    }
    wg.Wait()
	return result
}

func main() {
	// You can insert optional local tests here if desired.
}
