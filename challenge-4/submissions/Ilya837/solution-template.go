package main

import (
	"fmt"
	"sync"
)

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.

func ConcurrentBFS(graph map[int][]int, vInd int) []int {

	resMap := make(map[int]int, 0)

	next := make([]int, 0)

	next = append(next, vInd)

	resMap[vInd]++

	i := 0

	for {

		if i == len(next) {
			break
		}

		for _, v := range graph[next[i]] {
			if resMap[v] == 0 {
				resMap[v]++
				next = append(next, v)
			}
		}

		i++
	}

	return next
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {

	mute := sync.Mutex{}

	work := make(chan int)

	res := make(map[int][]int)

	go func() {
		for _, v := range queries {
			work <- v
		}

		close(work)

	}()

	wg := sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			for w := range work {
				workRes := ConcurrentBFS(graph, w)
				mute.Lock()
				res[w] = workRes
				mute.Unlock()
			}
			wg.Done()
		}()
	}

	wg.Wait()

	return res
}

func main() {
	graph := map[int][]int{
		0: {1, 2},
		1: {2, 3},
		2: {3},
		3: {4},
		4: {},
	}
	queries := []int{0, 1, 2}
	numWorkers := 2

	results := ConcurrentBFSQueries(graph, queries, numWorkers)
	/*
	   Possible output:
	   results[0] = [0 1 2 3 4]
	   results[1] = [1 2 3 4]
	   results[2] = [2 3 4]
	*/

	fmt.Println(results)
}
