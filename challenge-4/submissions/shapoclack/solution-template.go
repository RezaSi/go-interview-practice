package main

import (
    "sync"
    "fmt"
)


type BFSResult struct {
    StartNode int
    BFSOrder  []int
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
  
    if numWorkers <= 0 {
        return make(map[int][]int)
    }
    
    
    if len(queries) == 0 {
        return make(map[int][]int)
    }
    
    
    jobs := make(chan int, len(queries))
    
    
    results := make(chan BFSResult, len(queries))
    
    
    var wg sync.WaitGroup
    
    
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go bfsWorker(graph, jobs, results, &wg)
    }
    
    
    for _, startNode := range queries {
        jobs <- startNode
    }
    close(jobs) 
    
    
    go func() {
        wg.Wait()
        close(results) 
    }()
    
    
    finalResults := make(map[int][]int)
    for result := range results {
        finalResults[result.StartNode] = result.BFSOrder
    }
    
    return finalResults
}


func bfsWorker(graph map[int][]int, jobs <-chan int, results chan<- BFSResult, wg *sync.WaitGroup) {
    defer wg.Done()
    
   
    for startNode := range jobs {
        bfsOrder := performBFS(graph, startNode)
        results <- BFSResult{
            StartNode: startNode,
            BFSOrder:  bfsOrder,
        }
    }
}


func performBFS(graph map[int][]int, start int) []int {
    
    visited := make(map[int]bool)
    queue := []int{start}
    result := []int{}
    
    visited[start] = true
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        result = append(result, current)
        
       
        for _, neighbor := range graph[current] {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
    
    return result
}

func main() {
	
	graph := map[int][]int{
		1: {2, 3},
		2: {1, 4, 5},
		3: {1, 6},
		4: {2},
		5: {2, 6},
		6: {3, 5},
	}

	queries := []int{1, 2, 3, 4}

	results := ConcurrentBFSQueries(graph, queries, 2)

	for startNode, bfsOrder := range results {
		fmt.Printf("BFS from %d: %v\n", startNode, bfsOrder)
	}
}

