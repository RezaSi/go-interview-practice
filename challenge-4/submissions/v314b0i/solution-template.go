package main

import (
    "errors"
    "sync"
    "log"
)
// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.

var (
    ErrNodeNotInGraph = errors.New("Node not in Graph")
)

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
    results := map[int][]int{}
    if numWorkers < 1 {
        return results
    }
    resultsLock := sync.Mutex{}
    jobs := make(chan int, len(queries))
    
    // launch N workers
    wg := sync.WaitGroup{}
    for i := 0; i < numWorkers; i+=1 {
        wg.Add(1)
        go func(){
            defer wg.Done()
            for start := range jobs{
                result, err := BFS(graph, start)
                if err != nil {
                    log.Println(err)
                }
                resultsLock.Lock()
                results[start] = result
                resultsLock.Unlock()
            }
        }()
    }
    
    for _, start := range queries {
        jobs <- start
    }
    close(jobs)
    wg.Wait()
    return results
}

func BFS(graph map[int][]int, start int) (order []int, err error) {
    queued := map[int]bool{} // nodes alredy visited or added to queue
    queue := []int{} // current BFS queue
    
    // init
    queued[start] = true
    queue = append(queue, start)
    // BFS
    for len(queue) > 0 {
        // DeQueue next node from queue
        node := queue[0]
        order = append(order, node)
        queue = queue[1:]
        //if err = checkNodeInGraph(graph, node); err != nil {
        //    return
        //}
        // EnQueue node's neighbours, not already visited/queued
        for _, neighbour := range graph[node] {
            if _, exists := queued[neighbour]; !exists {
                queue = append(queue, neighbour)
                queued[neighbour] = true
            }
        }
    }
    return
}

func checkNodeInGraph(graph map[int][]int, start int) (err error) {
    if _, exists := graph[start]; !exists {
        err = ErrNodeNotInGraph
    }
    return
}

func main() {
	// You can insert optional local tests here if desired.
}
