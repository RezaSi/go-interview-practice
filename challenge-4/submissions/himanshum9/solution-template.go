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
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged
	
	output := make(map[int][]int)
// 	if len(graph) < 1 {
// 	    return output
// 	}
	var m sync.Mutex
	
	ch := make(chan int,len(queries))
	wg := sync.WaitGroup{}
	for i:=0;i < numWorkers;i++{
	    wg.Add(1)
	    go Worker(ch,&wg,output,graph,&m)
	}
	for _,value := range queries {
	    ch<-value
	}
	close(ch)
	wg.Wait()
	return output
}

func Worker(ch chan int,wg *sync.WaitGroup,output map[int][]int, graph map[int][]int,m *sync.Mutex) {
    defer wg.Done()
    for value := range ch {
	    visited := make([]int,len(graph) + 1)
	    innerQueue := []int{value}
	    res := []int{}
	    visited[value] = 1
	    
	    for len(innerQueue) != 0 {
	        firstEle := innerQueue[0]
	        res = append(res,firstEle)
	        innerQueue = innerQueue[1:len(innerQueue)]
	        visited[firstEle] = 1
	        val,ok := graph[firstEle]
    	    if ok {
    	        for _,j := range val {
    	            if visited[j] == 0 {
    	                visited[j] = 1
    	                innerQueue = append(innerQueue,j)
    	            }
    	        }
    	    }
	    }
	    m.Lock()
        output[value] = res
        m.Unlock()
	}
}

func main() {
	// You can insert optional local tests here if desired.
}
