package main

import (
	"sync"
)

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//

type startOrder struct {
	node  int
	order []int
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// startP channel sends queries to gorutines
	startP := make(chan int, len(queries))
	for _, s := range queries {
		startP <- s
	}
	close(startP)

	res := make(chan startOrder)
	var wg sync.WaitGroup

	// For loop makes workers as numWorkers param
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for s := range startP {
				order := bfsOrder(graph, s)
				res <- order
			}
		}()
	}

	// Separate gorunine wait for all workers finish ther jobs
	go func() {
		wg.Wait()
		close(res)
	}()

	// For loop reads res chanell antil it closed
	result := make(map[int][]int, len(queries))
	for r := range res {
		result[r.node] = r.order
	}
	return result
}

// bfsOrder makes queue and visited datastructure returns result for one query
func bfsOrder(g map[int][]int, s int) startOrder {
	visited := map[int]bool{s: true}
	queue := []int{s}
	order := startOrder{}
	order.node = s

	for len(queue) > 0 {
		curNode := queue[0]
		queue = queue[1:]
		order.order = append(order.order, curNode)

		for _, n := range g[curNode] {
			if !visited[n] {
				visited[n] = true
				queue = append(queue, n)
			}
		}
	}
	return order
}
func main() {
	// You can insert optional local tests here if desired.
}
