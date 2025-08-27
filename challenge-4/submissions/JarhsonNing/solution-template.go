package main

import (
	"sync"
)

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// 创建结果map和保护它的mutex
	results := make(map[int][]int)
	var mu sync.Mutex

	// 创建任务channel和WaitGroup
	taskChan := make(chan int, len(queries))
	var wg sync.WaitGroup

	// 启动worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for startNode := range taskChan {
				// 执行BFS
				bfsResult := bfs(graph, startNode)

				// 安全地将结果存储到map中
				mu.Lock()
				results[startNode] = bfsResult
				mu.Unlock()
			}
		}()
	}

	// 将所有查询发送到taskChan
	for _, query := range queries {
		taskChan <- query
	}

	// 关闭channel并等待所有workers完成
	close(taskChan)
	wg.Wait()

	return results
}

// bfs 执行广度优先搜索，返回遍历顺序
func bfs(graph map[int][]int, start int) []int {
	if _, exists := graph[start]; !exists {
		// 如果起始节点不在图中，返回包含该节点的单元素切片
		return []int{start}
	}

	visited := make(map[int]bool)
	queue := []int{start}
	result := []int{}

	for len(queue) > 0 {
		// 从队列前端取出节点
		current := queue[0]
		queue = queue[1:]

		// 如果已经访问过，跳过
		if visited[current] {
			continue
		}

		// 标记为已访问并添加到结果中
		visited[current] = true
		result = append(result, current)

		// 将邻接节点添加到队列中
		for _, neighbor := range graph[current] {
			if !visited[neighbor] {
				queue = append(queue, neighbor)
			}
		}
	}

	return result
}

func main() {
	// You can insert optional local tests here if desired.
}
