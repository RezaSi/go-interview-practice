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

type BFSResults struct {
        StartNode int
        Order     []int
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {

        result := map[int][]int{}

        if numWorkers == 0 {
                return result
        }

        results := make(chan BFSResults, len(queries))
        defer close(results)

        var wg sync.WaitGroup

        wg.Add(numWorkers)
        go pool(&wg, numWorkers, graph, queries, results)

        for i := 0; i < len(queries); i++ {
                res := <-results
                result[res.StartNode] = res.Order
        }

        wg.Wait()

        return result
}

func pool(wg *sync.WaitGroup, workers int, graph map[int][]int, tasks []int, results chan<- BFSResults) {
        jobs := make(chan int, len(tasks))
        for i := 0; i < workers; i++ {
                go worker(graph, jobs, results, wg)
        }

        for _, t := range tasks {
                jobs <- t
        }

        close(jobs)
}

func worker(graph map[int][]int, jobs <-chan int, results chan<- BFSResults, wg *sync.WaitGroup) {
        defer wg.Done()
        for j := range jobs {
                order := bfs(graph, j)
                results <- BFSResults{StartNode: j, Order: order}
        }
}

func bfs(graph map[int][]int, start int) []int {
        visited := map[int]bool{}
        queue := []int{start}
        order := []int{}

        visitedCounter := 0
        for {

                if len(graph) == visitedCounter {
                        break
                }

                tq := []int{}

                for _, q := range queue {
                        if visited[q] {
                                continue
                        }
                        visited[q] = true
                        visitedCounter++
                        tq = append(tq, graph[q]...)
                        order = append(order, q)
                }

                queue = append(queue, tq...)
                if len(queue) <= 1 {
                        break
                }
                queue = queue[1:]
                tq = []int{}
        }
        if len(graph) == 0 {
                return queue
        }
        return order
}



func main() {
	// You can insert optional local tests here if desired.
}
