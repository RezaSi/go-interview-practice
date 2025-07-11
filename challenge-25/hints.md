# Hints for Challenge 25: Graph Algorithms - Shortest Path

## Hint 1: BFS for Unweighted Graphs - Setup
Set up BFS with queue and distance tracking:
```go
func BreadthFirstSearch(graph [][]int, source int) ([]int, []int) {
    n := len(graph)
    distances := make([]int, n)
    predecessors := make([]int, n)
    visited := make([]bool, n)
    
    // Initialize distances to infinity
    for i := range distances {
        distances[i] = 1e9
        predecessors[i] = -1
    }
    
    distances[source] = 0
    
    // BFS queue implementation
    queue := []int{source}
    visited[source] = true
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        for _, neighbor := range graph[current] {
            if !visited[neighbor] {
                visited[neighbor] = true
                distances[neighbor] = distances[current] + 1
                predecessors[neighbor] = current
                queue = append(queue, neighbor)
            }
        }
    }
    
    return distances, predecessors
}
```

## Hint 2: Dijkstra's Algorithm - Priority Queue
Use a priority queue to always process the closest vertex first:
```go
import "container/heap"

type PriorityQueue []Node

type Node struct {
    vertex   int
    distance int
}

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].distance < pq[j].distance }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
    *pq = append(*pq, x.(Node))
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    node := old[n-1]
    *pq = old[0 : n-1]
    return node
}
```

## Hint 3: Dijkstra Implementation
Implement the main Dijkstra algorithm:
```go
func Dijkstra(graph [][]int, weights [][]int, source int) ([]int, []int) {
    n := len(graph)
    distances := make([]int, n)
    predecessors := make([]int, n)
    visited := make([]bool, n)
    
    for i := range distances {
        distances[i] = 1e9
        predecessors[i] = -1
    }
    
    distances[source] = 0
    
    pq := &PriorityQueue{}
    heap.Init(pq)
    heap.Push(pq, Node{vertex: source, distance: 0})
    
    for pq.Len() > 0 {
        current := heap.Pop(pq).(Node)
        vertex := current.vertex
        
        if visited[vertex] {
            continue
        }
        visited[vertex] = true
        
        for i, neighbor := range graph[vertex] {
            weight := weights[vertex][i]
            newDistance := distances[vertex] + weight
            
            if newDistance < distances[neighbor] {
                distances[neighbor] = newDistance
                predecessors[neighbor] = vertex
                heap.Push(pq, Node{vertex: neighbor, distance: newDistance})
            }
        }
    }
    
    return distances, predecessors
}
```

## Hint 4: Bellman-Ford Setup and Edge Relaxation
Implement edge relaxation for V-1 iterations:
```go
func BellmanFord(graph [][]int, weights [][]int, source int) ([]int, []bool, []int) {
    n := len(graph)
    distances := make([]int, n)
    predecessors := make([]int, n)
    hasPath := make([]bool, n)
    
    for i := range distances {
        distances[i] = 1e9
        predecessors[i] = -1
        hasPath[i] = false
    }
    
    distances[source] = 0
    hasPath[source] = true
    
    // Relax edges V-1 times
    for i := 0; i < n-1; i++ {
        for u := 0; u < n; u++ {
            if distances[u] == 1e9 {
                continue
            }
            
            for j, v := range graph[u] {
                weight := weights[u][j]
                if distances[u]+weight < distances[v] {
                    distances[v] = distances[u] + weight
                    predecessors[v] = u
                    hasPath[v] = true
                }
            }
        }
    }
    
    // Check for negative cycles...
    return distances, hasPath, predecessors
}
```

## Key Graph Algorithm Concepts:
- **BFS**: Level-by-level traversal for unweighted shortest paths
- **Dijkstra**: Greedy algorithm using priority queue for non-negative weights
- **Bellman-Ford**: Dynamic programming approach handling negative weights
- **Edge Relaxation**: Core operation of updating shortest distances
- **Priority Queue**: Efficient selection of minimum distance vertex
- **Negative Cycle Detection**: Essential for Bellman-Ford correctness 