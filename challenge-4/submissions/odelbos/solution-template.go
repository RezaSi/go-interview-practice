package main

import (
    "slices"
    "sync"
)

type bfsResult struct {
    start int
    result []int
}

func bfs(graph map[int][]int, start int) []int {
    var result []int
    queue := []int{start}
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        result = append(result, node)
        for _, child := range(graph[node]) {
            if ! slices.Contains(result, child) && ! slices.Contains(queue, child) {
                queue = append(queue, child)
            }
        }
    }
    return result
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
    if numWorkers == 0 {
        return map[int][]int{}
    }

    countQueries := len(queries)
    queriesChan := make(chan int, countQueries)
    resultsChan := make(chan bfsResult, countQueries)
    var wg sync.WaitGroup

    worker := func() {
        defer wg.Done()
        for query := range(queriesChan) {
            resultsChan <- bfsResult{query, bfs(graph, query)}
        }
    }

    // Create workers pool
    wg.Add(numWorkers)
    for i := 0; i < numWorkers; i++ {
        go worker()
    }

    // Send all queries to workers pool
    for _, q := range queries {
        queriesChan <- q
    }
    close(queriesChan)

    // Wait ...
    go func() {
        wg.Wait()
        close(resultsChan)
    }()

    // Collect results
    results := make(map[int][]int)
    for res := range(resultsChan) {
        results[res.start] = res.result
    }

    return results
}

func main() {
}
