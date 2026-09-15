package main

import (
	"container/heap"
	"errors"
	"fmt"
)

// main demonstrates the three graph-traversal algorithms implemented in this
// file (BFS, Dijkstra, Bellman-Ford) on small example graphs and prints the
// resulting distances and predecessor arrays.
func main() {
	// Example 1: Unweighted graph for BFS.
	unweightedGraph := [][]int{
		{1, 2},    // Vertex 0 has edges to vertices 1 and 2
		{0, 3, 4}, // Vertex 1 has edges to vertices 0, 3, and 4
		{0, 5},    // Vertex 2 has edges to vertices 0 and 5
		{1},       // Vertex 3 has an edge to vertex 1
		{1},       // Vertex 4 has an edge to vertex 1
		{2},       // Vertex 5 has an edge to vertex 2
	}

	// Run BFS from vertex 0 and print the resulting arrays.
	distances, predecessors := BreadthFirstSearch(unweightedGraph, 0)
	fmt.Println("BFS Results:")
	fmt.Printf("Distances: %v\n", distances)
	fmt.Printf("Predecessors: %v\n", predecessors)
	fmt.Println()

	// Example 2: Weighted graph for Dijkstra (non-negative weights).
	weightedGraph := [][]int{
		{1, 2},    // Vertex 0 has edges to vertices 1 and 2
		{0, 3, 4}, // Vertex 1 has edges to vertices 0, 3, and 4
		{0, 5},    // Vertex 2 has edges to vertices 0 and 5
		{1},       // Vertex 3 has an edge to vertex 1
		{1},       // Vertex 4 has an edge to vertex 1
		{2},       // Vertex 5 has an edge to vertex 2
	}
	weights := [][]int{
		{5, 10},   // Edge from 0 to 1 has weight 5, edge from 0 to 2 has weight 10
		{5, 3, 2}, // Edge weights from vertex 1
		{10, 2},   // Edge weights from vertex 2
		{3},       // Edge weights from vertex 3
		{2},       // Edge weights from vertex 4
		{2},       // Edge weights from vertex 5
	}

	// Run Dijkstra from vertex 0 and print the resulting arrays.
	dijkstraDistances, dijkstraPredecessors := Dijkstra(weightedGraph, weights, 0)
	fmt.Println("Dijkstra Results:")
	fmt.Printf("Distances: %v\n", dijkstraDistances)
	fmt.Printf("Predecessors: %v\n", dijkstraPredecessors)
	fmt.Println()

	// Example 3: Directed graph with a negative-weight edge for Bellman-Ford.
	negativeWeightGraph := [][]int{
		{1, 2},
		{3},
		{1, 3},
		{4},
		{},
	}
	negativeWeights := [][]int{
		{6, 7},  // Edge weights from vertex 0
		{5},     // Edge weights from vertex 1
		{-2, 4}, // Edge weights from vertex 2 (note the negative weight on 2 -> 1)
		{2},     // Edge weights from vertex 3
		{},      // Edge weights from vertex 4
	}

	// Run Bellman-Ford from vertex 0 and print distances, hasPath, and predecessors.
	bfDistances, hasPath, bfPredecessors := BellmanFord(negativeWeightGraph, negativeWeights, 0)
	fmt.Println("Bellman-Ford Results:")
	fmt.Printf("Distances: %v\n", bfDistances)
	fmt.Printf("Has Path: %v\n", hasPath)
	fmt.Printf("Predecessors: %v\n", bfPredecessors)
}

// Queue is a minimal FIFO queue of ints used by BreadthFirstSearch.
type Queue []int

// Push appends a value to the tail of the queue.
// The error return is unused but kept to mirror a typical queue interface.
func (q *Queue) Push(v int) error {
	*q = append(*q, v)
	return nil
}

// Pop removes and returns the value at the head of the queue.
// It returns an error if the queue is empty.
func (q *Queue) Pop() (int, error) {
	if len(*q) == 0 {
		return 0, errors.New("empty queue")
	}
	oldQ := *q
	v := oldQ[0]
	*q = oldQ[1:]
	return v, nil
}

// BreadthFirstSearch implements BFS for unweighted graphs to find shortest paths
// from a source vertex to all other vertices.
// Returns:
//   - distances: slice where distances[i] is the shortest distance from source to vertex i
//     (1e9 sentinel value if vertex i is unreachable).
//   - predecessors: slice where predecessors[i] is the vertex that comes before i in the
//     shortest path (-1 if i is the source or unreachable).
func BreadthFirstSearch(graph [][]int, source int) ([]int, []int) {
	l := len(graph)

	// Initialize distances to the "infinity" sentinel and predecessors to "none".
	predecessors := make([]int, l)
	distances := make([]int, l)
	for i, _ := range distances {
		distances[i] = 1000000000
		predecessors[i] = -1
	}

	distances[source] = 0

	queue := Queue{}

	// visited prevents re-enqueueing of vertices already discovered.
	visited := make(map[int]struct{})

	visited[source] = struct{}{}

	queue.Push(source)
	for len(queue) > 0 {
		// Dequeue the next vertex and relax its outgoing edges.
		v, _ := queue.Pop()

		for _, nextV := range graph[v] {
			if _, exists := visited[nextV]; exists {
				continue
			}
			visited[nextV] = struct{}{}
			distance := distances[v] + 1
			queue.Push(nextV)
			// In BFS the first discovery is always the shortest path,
			// so this guard mainly protects against re-assignment.
			if distances[nextV] > distance {
				distances[nextV] = distance
				predecessors[nextV] = v
			}
		}
	}
	return distances, predecessors
}

// ===============================
// Priority-queue implementation backed by container/heap, used by Dijkstra.

// PriorityQueue is a min-heap of *Vertex ordered by Vertex.distance.
type PriorityQueue []*Vertex

// Vertex represents a node in the priority queue together with its current
// best-known distance from the source.
type Vertex struct {
	graphIndex int // index of the vertex in the original graph slice
	distance   int // current best-known distance from source
	index      int // index within the heap (maintained by heap.Interface)
}

// NewPriorityVertex constructs a *Vertex with the given graph index and distance.
func NewPriorityVertex(graphIndex, distance int) *Vertex {
	return &Vertex{
		graphIndex: graphIndex,
		distance:   distance,
	}
}

// Len returns the number of elements in the queue. Required by heap.Interface.
func (q PriorityQueue) Len() int {
	return len(q)
}

// Less compares two vertices by distance, enforcing a min-heap order.
func (q PriorityQueue) Less(i, j int) bool {
	return q[i].distance < q[j].distance
}

// Swap exchanges two vertices in the heap and updates their stored heap indices.
func (q PriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

// Push appends a new vertex to the underlying slice. Required by heap.Interface;
// it should be invoked through heap.Push, not directly.
func (q *PriorityQueue) Push(x interface{}) {
	n := len(*q)
	vertex := x.(*Vertex)
	vertex.index = n
	*q = append(*q, vertex)
}

// Pop removes and returns the last element of the underlying slice.
// heap.Pop maintains the heap invariant before calling this method.
func (q *PriorityQueue) Pop() interface{} {
	n := len(*q)
	c := *q
	v := c[n-1]
	c[n-1] = nil
	*q = c[:n-1]
	return v
}

// Dijkstra implements Dijkstra's algorithm for weighted graphs with non-negative
// weights to find shortest paths from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func Dijkstra(graph [][]int, weights [][]int, source int) ([]int, []int) {
	l := len(graph)

	// Initialize distances to the "infinity" sentinel and predecessors to "none".
	predecessors := make([]int, l)
	distances := make([]int, l)
	for i, _ := range distances {
		distances[i] = 1000000000
		predecessors[i] = -1
	}

	distances[source] = 0

	queue := PriorityQueue{}

	// visited marks vertices whose shortest distance has been finalised (popped).
	visited := make(map[int]struct{})

	// Seed the priority queue with the source at distance 0.
	heap.Init(&queue)
	heap.Push(&queue, NewPriorityVertex(source, 0))

	for len(queue) > 0 {
		// Pop the vertex with the smallest tentative distance.
		x := heap.Pop(&queue)
		vertex := x.(*Vertex)

		v := vertex.graphIndex
		//Finalize v only when popped; skip state duplicate heap entries.
		if _, done := visited[v]; done {
			continue
		}
		visited[v] = struct{}{}

		// Relax each outgoing edge of v.
		for i, nextV := range graph[v] {
			if _, done := visited[nextV]; done {
				continue
			}
			distance := distances[v] + weights[v][i]

			if distances[nextV] > distance {
				distances[nextV] = distance
				predecessors[nextV] = v
				heap.Push(&queue, NewPriorityVertex(nextV, distance))
			}
		}
	}
	return distances, predecessors
}

// BellmanFord implements the Bellman-Ford algorithm for weighted graphs that may
// contain negative weight edges to find shortest paths from a source vertex to
// all other vertices. It also detects which vertices have a well-defined shortest
// path (i.e. are reachable from source AND are not influenced by a negative cycle).
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - hasPath: slice where hasPath[i] is true if there is a path from source to i without a negative cycle
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func BellmanFord(graph [][]int, weights [][]int, source int) ([]int, []bool, []int) {
	l := len(graph)

	// Initialize distances to the "infinity" sentinel and predecessors to "none".
	predecessors := make([]int, l)
	distances := make([]int, l)
	hasPath := make([]bool, l)
	for i := range distances {
		distances[i] = 1000000000
		predecessors[i] = -1
	}

	distances[source] = 0

	// Standard Bellman-Ford: relax every edge V-1 times. After V-1 rounds all
	// shortest paths that do not pass through a negative cycle are finalised.
	for k := 1; k < l; k++ {
		for i := 0; i < l; i++ {
			for j := 0; j < len(graph[i]); j++ {
				s := i
				d := graph[i][j]
				w := weights[i][j]
				// Skip relaxation from an unreachable vertex to avoid spurious updates.
				if distances[i] != 1000000000 && distances[s]+w < distances[d] {
					distances[d] = distances[s] + w
					predecessors[d] = s
				}
			}
		}
	}

	// A vertex has a path iff its distance is finite after V-1 iterations.
	for i := 0; i < l; i++ {
		hasPath[i] = distances[i] != 1000000000
	}

	// Negative-cycle detection: any edge that can still be relaxed indicates
	// that its destination is reachable from a negative-weight cycle.
	for i := 0; i < l; i++ {
		for j := 0; j < len(graph[i]); j++ {
			s := i
			d := graph[i][j]
			w := weights[i][j]
			if distances[i] != 1000000000 && distances[s]+w < distances[d] {
				hasPath[d] = false
			}
		}
	}

	// Propagate the "no valid shortest path" flag to all descendants of tainted
	// vertices: any vertex reachable from a negative cycle can have its distance
	// driven arbitrarily low, so its shortest path is undefined as well.
	for k := 1; k < l; k++ {
		for i := 0; i < l; i++ {
			for j := 0; j < len(graph[i]); j++ {
				if !hasPath[i] {
					hasPath[graph[i][j]] = false
				}
			}
		}
	}

	return distances, hasPath, predecessors
}
