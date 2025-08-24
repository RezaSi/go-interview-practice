package main

import (
	"slices"
)

func bfs(graph map[int][]int, s int) []int {
	var visited []int
	var queue []int

	visited = append(visited, s)
	queue = append(queue, s)

	var v int // current vertex
	for len(queue) > 0 {
		v, queue = queue[0], queue[1:]
		neighbours := graph[v]
		for _, n := range neighbours {
			if !slices.Contains(visited, n) {
				visited = append(visited, n)
				queue = append(queue, n)
			}
		}
	}

	return visited
}

func worker(graph map[int][]int, jobs <-chan int, results chan<- []int) {
	for j := range jobs {
		results <- bfs(graph, j)
	}
}

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	numJobs := len(queries)
	jobs := make(chan int, numJobs)
	results := make(chan []int, numJobs)

	if numWorkers == 0 {
		return map[int][]int{}
	}

	for w := 1; w <= numWorkers; w++ {
		go worker(graph, jobs, results)
	}

	for i := 1; i <= numJobs; i++ {
		jobs <- queries[i-1]
	}
	close(jobs)

	answer := make(map[int][]int)
	for i := 1; i <= numJobs; i++ {
		result := <-results
		answer[result[0]] = result
	}

	return answer
}

func main() {
	// You can insert optional local tests here if desired.
	graph := map[int][]int{
		0: {1, 2},
		1: {2, 3},
		2: {3},
		3: {4},
		4: {},
	}
	queries := []int{0, 1, 2}
	numWorkers := 2

	ConcurrentBFSQueries(graph, queries, numWorkers)
}
