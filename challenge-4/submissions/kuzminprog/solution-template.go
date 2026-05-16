package main

import (
	"sync"
)

type RingQueue[T any] struct {
	buf  []T
	head int
	tail int
	size int
}

func NewRingQueue[T any](cap int) *RingQueue[T] {
	if cap < 1 {
		cap = 1
	}
	return &RingQueue[T]{buf: make([]T, cap)}
}

func (q *RingQueue[T]) Push(v T) {
	if q.size == len(q.buf) {
		q.grow()
	}
	q.buf[q.tail] = v
	q.tail = (q.tail + 1) % len(q.buf)
	q.size++
}

func (q *RingQueue[T]) Pop() T {
	var zero T
	if q.size == 0 {
		return zero
	}
	v := q.buf[q.head]
	q.buf[q.head] = zero
	q.head = (q.head + 1) % len(q.buf)
	q.size--
	return v
}

func (q *RingQueue[T]) grow() {
	nb := make([]T, len(q.buf)*2)
	for i := 0; i < q.size; i++ {
		nb[i] = q.buf[(q.head+i)%len(q.buf)]
	}
	q.buf = nb
	q.head = 0
	q.tail = q.size
}

func BFSQueries(graph map[int][]int, start int) []int {
	visited := make(map[int]struct{})
	queue := NewRingQueue[int](len(graph))

	visited[start] = struct{}{}
	queue.Push(start)

	order := make([]int, 0)
	for queue.size > 0 {
		v := queue.Pop()
		order = append(order, v)

		for _, u := range graph[v] {
			if _, ok := visited[u]; !ok {
				visited[u] = struct{}{}
				queue.Push(u)
			}
		}
	}

	return order
}

type BFSJob struct {
	graph map[int][]int
	start int
}

type BFSResult struct {
	order []int
	start int
}

func BFSWorker(jobs <-chan BFSJob, results chan<- BFSResult) {
	for j := range jobs {
		order := BFSQueries(j.graph, j.start)
		results <- BFSResult{order: order, start: j.start}
	}
}

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	jobs := make(chan BFSJob, len(queries))
	results := make(chan BFSResult, len(queries))

	var wg sync.WaitGroup

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			BFSWorker(jobs, results)
		}()
	}

	go func() {
		for _, q := range queries {
			jobs <- BFSJob{graph: graph, start: q}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	orders := make(map[int][]int, len(queries))
	for res := range results {
		orders[res.start] = res.order
	}

	return orders
}

func main() {
	// You can insert optional local tests here if desired.
}
