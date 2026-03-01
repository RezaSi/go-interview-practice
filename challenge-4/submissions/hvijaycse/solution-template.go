package main

import "sync"

type Result struct {
	start int
	order []int
}

// standard BFS algo.
func bfs(graph map[int][]int, start int) []int {

	order := []int{}
	visited := map[int]bool{}
	queue := []int{start}

	visited[start] = true

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		order = append(order, node)

		for _, neighbor := range graph[node] {

			if visited[neighbor] {
				continue
			}
			visited[neighbor] = true
			queue = append(queue, neighbor)

		}
	}

	return order
}

// worker to perform  bfs.
func worker(wg *sync.WaitGroup, graph map[int][]int, jobs <-chan int, result chan<- Result) {

	defer wg.Done()

	for job := range jobs {
		order := bfs(graph, job)
		result <- Result{start: job, order: order}
	}

}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {

	numWorkers = min(numWorkers, len(queries))
	output := make(map[int][]int)

	if numWorkers <= 0 {
		return output
	}
	wg := &sync.WaitGroup{}
	jobs := make(chan int)
	results := make(chan Result)

	wg.Add(numWorkers)

	// Creating workers to execute task.
	for range numWorkers {
		go worker(wg, graph, jobs, results)
	}

	// Passing the queries to workers using jobs channel
	go func() {
		for _, query := range queries {
			jobs <- query
		}
		// closing jobs since all the
		// queries are passed
		close(jobs)
	}()

	// Close results AFTER workers finish
	go func() {
		wg.Wait()
		// Asap the workers are done
		// we close the results channel
		// otherwise the range will keep looping
		// on it forever.
		close(results)
	}()

	for result := range results {
		output[result.start] = result.order
	}

	return output
}

func main() {
	// You can insert optional local tests here if desired.
}
