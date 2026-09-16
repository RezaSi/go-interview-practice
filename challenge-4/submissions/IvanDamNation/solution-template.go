package main

import "sync"

type result struct {
    query   int
    path    []int
}

type bfsWorkerBuf struct {
    path    []int
    visited map[int]struct{}
}

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	paths := make(map[int][]int, len(queries))
	
	if len(queries) == 0 || numWorkers <= 0 {
	    return paths
	}
	
	qCh := make(chan int, len(queries))
	resCh := make(chan result, len(queries))
	
	if len(queries) < numWorkers {
	    numWorkers = len(queries)
	}
	
	maxNodes := len(graph)
	
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for range numWorkers {
	    go func() {
	        defer wg.Done()
	        
	        buf := &bfsWorkerBuf{
	            path: make([]int, 0, maxNodes),
	            visited: make(map[int]struct{}, maxNodes),
	        }
	        
	        for q := range qCh {
	            p := bfsWorker(graph, q, buf)
	            
	            resCh <-result{
	                query: q,
	                path: p,
	            }
	        }
	    }()
	}
	
	for _, query := range queries {
	    qCh <-query
	}
	close(qCh)
	
	go func() {
	    wg.Wait()
	    close(resCh)
	}()
	
	for res := range resCh {
	    paths[res.query] = res.path
	}
	
	return paths
}

func bfsWorker(graph map[int][]int, start int, buf *bfsWorkerBuf) []int {
    buf.path = buf.path[:0]
    clear(buf.visited)
    
    buf.path = append(buf.path, start)
    buf.visited[start] = struct{}{}
    
    ptr := 0
    for ptr < len(buf.path) {
        node := buf.path[ptr]
        ptr++
        
        for _, v := range graph[node] {
            if _, exists := buf.visited[v]; exists {
                continue
            }
            
            buf.visited[v] = struct{}{}
            buf.path = append(buf.path, v)
        }
    }
    res := make([]int, len(buf.path))
    copy(res, buf.path)
    
    return res
}

func main() {
	// You can insert optional local tests here if desired.
}
