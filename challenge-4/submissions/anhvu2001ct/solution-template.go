package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func concurrentlyRun[T any](jobs []T, numWorkers int, workFn func(idx int, job T)) {
	if numWorkers == 0 {
		return
	}

	counter := atomic.Int32{}
	n := len(jobs)
	wg := sync.WaitGroup{}

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				i := int(counter.Add(1) - 1)
				if i >= n {
					break
				}
				workFn(i, jobs[i])
			}
		}()
	}

	wg.Wait()
}

// start bfs at node u
func bfs(graph map[int][]int, u int) []int {
	queue := []int{}
	marked := map[int]struct{}{}
	result := []int{}

	queue = append(queue, u)
	marked[u] = struct{}{}

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		result = append(result, u)

		for _, v := range graph[u] {
			if _, exist := marked[v]; exist {
				continue
			}

			marked[v] = struct{}{}
			queue = append(queue, v)
		}
	}

	return result
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	if numWorkers == 0 {
		return nil
	}

	ans := map[int][]int{}
	result := make([][]int, len(queries))

	concurrentlyRun(queries, numWorkers, func(idx, query int) {
		result[idx] = bfs(graph, query)
	})

	for i, res := range result {
		ans[queries[i]] = res
	}

	return ans
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
	/*
	   Possible output:
	   results[0] = [0 1 2 3 4]
	   results[1] = [1 2 3 4]
	   results[2] = [2 3 4]
	*/
}
