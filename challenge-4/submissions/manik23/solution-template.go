package main

import (
	"container/list"
	"sync"
)


type reachableNode struct {
	node  int
	neigh []int
}

func process(graph map[int][]int, jobs <-chan int, resultBFS chan reachableNode) {
	for startNode := range jobs {
		reachableNode := travel(graph, startNode)
		resultBFS <- reachableNode
	}
}

func travel(graph map[int][]int, startNode int) reachableNode {

	visited := make(map[int]bool)

	list := list.New()
	var order []int

	list.PushBack(startNode)
	visited[startNode] = true

	for list.Len() > 0 {

		node := list.Front().Value.(int)
		list.Remove(list.Front())
		order = append(order, node)

		for _, neigh := range graph[node] {

			if !visited[neigh] {
				visited[neigh] = true
				list.PushBack(neigh)

			}
		}
	}

	return reachableNode{
		node:  startNode,
		neigh: order,
	}
}

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {

	if numWorkers == 0 || len(queries) == 0 {

		return map[int][]int{}
	}

	if len(graph) == 0 {
		finalResult := make(map[int][]int)
		for _, startNode := range queries {
			finalResult[startNode] = []int{startNode}
		}
		return finalResult
	}

	jobs := make(chan int, len(queries))

	resultBFS := make(chan reachableNode, len(queries))
	defer close(resultBFS)

	var wg sync.WaitGroup

	for range numWorkers {

		wg.Add(1)
		go func() {
			defer wg.Done()
			process(graph, jobs, resultBFS)
		}()

	}

	go func() {
		for _, startNode := range queries {
			jobs <- startNode
		}
		close(jobs)
	}()

	wg.Wait()

	resultMap := make(map[int][]int)
	for range queries {
		result := <-resultBFS
		resultMap[result.node] = result.neigh
	}

	return resultMap
}

func main() {
	// You can insert optional local tests here if desired.
}
