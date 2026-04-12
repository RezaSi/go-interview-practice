package main

import "sync"

func BFS(source int, graph map[int][]int) []int {
	mark := make([]bool, len(graph))
	index := 0
	queue := make([]int, 0)
	var u int

	queue = append(queue, source)
	mark[source] = true
	for index < len(queue) {
		u = queue[index]
		index++
		for _, v := range graph[u] {
			if !mark[v] {
				mark[v] = true
				queue = append(queue, v)
			}
		}
	}
	return queue
}

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	answer := make(map[int][]int)

	if numWorkers == 0 {
		return answer
	} else if len(graph) == 0 {
		graph[0] = []int{}
	}

	c := make(chan interface{}, numWorkers)
	for i := 0; i < numWorkers; i++ {
		c <- interface{}(true)
	}

	var wg sync.WaitGroup
	var mut sync.Mutex
	wg.Add(len(queries))
	for _, source := range queries {
		<-c
		go func(s int) {
			mut.Lock()
			answer[s] = BFS(s, graph)
			mut.Unlock()
			c <- interface{}(true)
			wg.Done()
		}(source)
	}
	wg.Wait()

	return answer
}

func main() {
	// You can insert optional local tests here if desired.
}
