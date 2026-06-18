package main

import (
	"fmt"
	"sync"
)

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// graph[u] = []int{v1, v2, ...}
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	results := sync.Map{}
	jobs := make(chan int, len(queries))

	// 投递所有查询任务
	for _, q := range queries {
		jobs <- q
	}
	close(jobs)

	var wg sync.WaitGroup

	// 启动固定数量的 worker
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for start := range jobs {
				visited := make(map[int]bool)
				order := bfs(graph, start, visited)
				results.Store(start, order)
			}
		}()
	}

	wg.Wait()

	// 转成普通 map 返回
	output := make(map[int][]int)
	results.Range(func(key, value any) bool {
		output[key.(int)] = value.([]int)
		return true
	})
	return output
}

// 单个 BFS 实现
func bfs(graph map[int][]int, start int, visited map[int]bool) []int {
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

		for _, nei := range graph[node] {
			if !visited[nei] {
				queue = append(queue, nei)
			}
		}
	}
	return order
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
