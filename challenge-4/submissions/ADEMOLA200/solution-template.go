package main

import (
	"sync"
	"fmt"
)

// bfs performs a standard BFS traversal from a start node
func bfs(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	order := []int{}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if visited[node] {
			continue
		}
		visited[node] = true
		order = append(order, node)

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				queue = append(queue, neighbor)
			}
		}
	}
	return order
}

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	type job struct {
		start int
	}
	type result struct {
		start int
		order []int
	}

	jobs := make(chan job)
	results := make(chan result)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				order := bfs(graph, j.start)
				results <- result{start: j.start, order: order}
			}
		}()
	}

	// Send jobs
	go func() {
		for _, q := range queries {
			jobs <- job{start: q}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	output := make(map[int][]int)
	for res := range results {
		output[res.start] = res.order
	}

	return output
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
	for start, order := range results {
		println("Start:", start, "→", fmtSlice(order))
	}
}

func fmtSlice(s []int) string {
	str := "["
	for i, v := range s {
		if i > 0 {
			str += " "
		}
		str += fmt.Sprint(v)
	}
	str += "]"
	return str
}
