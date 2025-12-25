package main

import (
	"fmt"
)

func BFS(graph map[int][]int, start int) []int {
	seen := map[int]bool{}
	seen[start] = true
	left := 0
	q := []int{start}
	for left < len(q) {
		u := q[left]
		left++
		for _, v := range graph[u] {
			if seen[v] == false {
				q = append(q, v)
				seen[v] = true
			}
		}

	}
	return q
}

// func worker(graph map[int][]int, jobs <-chan int, results chan<- BFSResult, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for start := range jobs {
// 		order := BFS(graph, start)
// 		results <- BFSResult{StartNode: start, Order: order}
// 	}
// }

type BFSResult struct {
	StartNode int
	Order     []int
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	if len(queries) == 0 || numWorkers == 0 {
		return make(map[int][]int)
	}
	jobs := make(chan int, len(queries))
	results := make(chan BFSResult, len(queries))
	// var wg sync.WaitGroup  //usefull when we dont know the number of queries

	for range numWorkers {
		// wg.Add(1)
		// go worker(graph, jobs, results, &wg)
		go func() {
			// defer wg.Done()
			for start := range jobs {
				order := BFS(graph, start)
				go func() {
					results <- BFSResult{StartNode: start, Order: order}
				}()
			}
		}()
	}

	go func() {
		for _, x := range queries {
			jobs <- x
		}
		close(jobs)
	}()

	// go func() {
	// 	// wg.Wait()

	// }()
	defer close(results)
	resultMap := make(map[int][]int)
	for range len(queries) {
		result := <-results
		resultMap[result.StartNode] = result.Order
	}

	return resultMap
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
	fmt.Println(results)
}
