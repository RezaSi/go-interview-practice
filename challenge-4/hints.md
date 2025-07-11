# Hints for Concurrent Graph BFS Queries

## Hint 1: Worker Pool Pattern
You need to implement a worker pool pattern. Create `numWorkers` goroutines that will process BFS queries from a channel:
```go
jobs := make(chan int, len(queries))
results := make(chan BFSResult, len(queries))
```

## Hint 2: BFS Result Structure
Define a struct to hold BFS results:
```go
type BFSResult struct {
    StartNode int
    Order     []int
}
```

## Hint 3: Starting Workers
Launch the specified number of worker goroutines:
```go
for i := 0; i < numWorkers; i++ {
    go worker(graph, jobs, results)
}
```

## Hint 4: Sending Jobs
Send all queries to the jobs channel and close it:
```go
for _, query := range queries {
    jobs <- query
}
close(jobs)
```

## Hint 5: Standard BFS Implementation
Each worker performs standard BFS using a queue:
```go
func bfs(graph map[int][]int, start int) []int {
    visited := make(map[int]bool)
    queue := []int{start}
    order := []int{}
    // ... implement BFS
}
```

## Hint 6: Worker Function Structure
The worker function should process jobs until the channel is closed:
```go
func worker(graph map[int][]int, jobs <-chan int, results chan<- BFSResult) {
    for start := range jobs {
        order := bfs(graph, start)
        results <- BFSResult{StartNode: start, Order: order}
    }
}
```

## Hint 7: Collecting Results
Collect all results and convert to the required map format:
```go
resultMap := make(map[int][]int)
for i := 0; i < len(queries); i++ {
    result := <-results
    resultMap[result.StartNode] = result.Order
}
```

## Hint 8: BFS Queue Operations
For BFS, use slice operations for queue:
```go
// Dequeue (remove from front)
node := queue[0]
queue = queue[1:]

// Enqueue (add to back)
queue = append(queue, neighbor)
``` 