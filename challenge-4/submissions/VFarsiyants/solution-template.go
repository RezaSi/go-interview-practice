package main

type BFSResult struct {
	StartNode int
	Order     []int
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	jobs := make(chan int, len(queries))
	results := make(chan BFSResult, len(queries))
	
	if numWorkers == 0 {
	   return map[int][]int{}
	}

	for range numWorkers {
		go worker(graph, jobs, results)
	}

	for _, query := range queries {
		jobs <- query
	}
	close(jobs)

	resultMap := make(map[int][]int)
	for range queries {
		result := <-results
		resultMap[result.StartNode] = result.Order
	}
	return resultMap
}

func worker(graph map[int][]int, jobs chan int, results chan BFSResult) {
	for start := range jobs {
		order := bfs(graph, start)
		results <- BFSResult{StartNode: start, Order: order}
	}
}

func bfs(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	queue := []int{start}
	order := []int{}

	visited[start] = true

	for len(queue) != 0 {
		currentNode := queue[0]
		queue = queue[1:]
		childNodes := graph[currentNode]
		order = append(order, currentNode)

		for _, childNode := range childNodes {
			if !visited[childNode] {
				visited[childNode] = true
				queue = append(queue, childNode)
			}
		}
	}

	return order
}

func main() {
	// You can insert optional local tests here if desired.
}
