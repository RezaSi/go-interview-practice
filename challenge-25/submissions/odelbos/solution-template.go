package main

import (
	"fmt"
	"container/heap"
)

func main() {
	// Example 1: Unweighted graph for BFS
	unweightedGraph := [][]int{
		{1, 2},    // Vertex 0 has edges to vertices 1 and 2
		{0, 3, 4}, // Vertex 1 has edges to vertices 0, 3, and 4
		{0, 5},    // Vertex 2 has edges to vertices 0 and 5
		{1},       // Vertex 3 has an edge to vertex 1
		{1},       // Vertex 4 has an edge to vertex 1
		{2},       // Vertex 5 has an edge to vertex 2
	}

	// Test BFS
	distances, predecessors := BreadthFirstSearch(unweightedGraph, 0)
	fmt.Println("BFS Results:")
	fmt.Printf("Distances: %v\n", distances)
	fmt.Printf("Predecessors: %v\n", predecessors)
	fmt.Println()

	// Example 2: Weighted graph for Dijkstra
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

	// Test Dijkstra
	dijkstraDistances, dijkstraPredecessors := Dijkstra(weightedGraph, weights, 0)
	fmt.Println("Dijkstra Results:")
	fmt.Printf("Distances: %v\n", dijkstraDistances)
	fmt.Printf("Predecessors: %v\n", dijkstraPredecessors)
	fmt.Println()

	// Example 3: Graph with negative weights for Bellman-Ford
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
		{-2, 4}, // Edge weights from vertex 2 (note the negative weight)
		{2},     // Edge weights from vertex 3
		{},      // Edge weights from vertex 4
	}

	// Test Bellman-Ford
	bfDistances, hasPath, bfPredecessors := BellmanFord(negativeWeightGraph, negativeWeights, 0)
	fmt.Println("Bellman-Ford Results:")
	fmt.Printf("Distances: %v\n", bfDistances)
	fmt.Printf("Has Path: %v\n", hasPath)
	fmt.Printf("Predecessors: %v\n", bfPredecessors)
}

// -----------------------------------------------------------
// BreadthFirstSearch
// -----------------------------------------------------------

// BreadthFirstSearch implements BFS for unweighted graphs to find shortest paths
// from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func BreadthFirstSearch(graph [][]int, source int) ([]int, []int) {
	nbVertices := len(graph)
	distances := make([]int, nbVertices)
	predecessors := make([]int, nbVertices)
	queue := []int{}
	for i := range(nbVertices) {
		distances[i] = int(1e9)
		predecessors[i] = -1
	}

	distances[source] = 0
	queue = append(queue, source)
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range(graph[u]) {
			if distances[v] == int(1e9) {
				distances[v] = distances[u] + 1
				predecessors[v] = u
				queue = append(queue, v)
			}
		}
	}
	return distances, predecessors
}

// -----------------------------------------------------------
// Dijkstra
// -----------------------------------------------------------

type Item struct {
	vertex   int
	distance int
}

type PriorityQueue []*Item

// heap.Interface
func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(*Item))
}

func (pq *PriorityQueue) Pop() any {
	n := len(*pq)
	item := (*pq)[n - 1]
	*pq = (*pq)[:n - 1]
	return item
}

// Dijkstra implements Dijkstra's algorithm for weighted graphs with non-negative weights
// to find shortest paths from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func Dijkstra(graph [][]int, weights [][]int, source int) ([]int, []int) {
	nbVertices := len(graph)
	distances := make([]int, nbVertices)
	predecessors := make([]int, nbVertices)
	pq := make(PriorityQueue, 0)
	for i := range(nbVertices) {
		distances[i] = int(1e9)
		predecessors[i] = -1
	}

	distances[source] = 0
	heap.Push(&pq, &Item{vertex: source, distance: 0})

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		u := item.vertex

		if item.distance > distances[u] {
			continue
		}

		for i, v := range(graph[u]) {
			weight := weights[u][i]
			if distances[u] + weight < distances[v] {
				distances[v] = distances[u] + weight
				predecessors[v] = u
				heap.Push(&pq, &Item{vertex: v, distance: distances[v]})
			}
		}
	}

	return distances, predecessors
}

// -----------------------------------------------------------
// BellmanFord
// -----------------------------------------------------------

// BellmanFord implements the Bellman-Ford algorithm for weighted graphs that may contain
// negative weight edges to find shortest paths from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - hasPath: slice where hasPath[i] is true if there is a path from source to i without a negative cycle
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func BellmanFord(graph [][]int, weights [][]int, source int) ([]int, []bool, []int) {
	nbVertices := len(graph)
	distances := make([]int, nbVertices)
	predecessors := make([]int, nbVertices)
	for i := range distances {
		distances[i] = int(1e9)
		predecessors[i] = -1
	}
	distances[source] = 0

	for range(nbVertices - 1) {
		for v := range(nbVertices) {
			for j, u := range graph[v] {
				w := weights[v][j]
				if distances[v] != int(1e9) && distances[v] + w < distances[u] {
					distances[u] = distances[v] + w
					predecessors[u] = v
				}
			}
		}
	}

	hasPath := make([]bool, nbVertices)
	for i := range hasPath {
		hasPath[i] = distances[i] != int(1e9)
	}
	negCycle := make([]bool, nbVertices)

	for v := range(nbVertices) {
		for j, u := range(graph[v]) {
			w := weights[v][j]
			if distances[v] != int(1e9) && distances[v] + w < distances[u] {
				negCycle[v] = true
			}
		}
	}

	changed := true
	for changed {
		changed = false
		for v := range(nbVertices) {
			if negCycle[v] {
				for _, u := range(graph[v]) {
					if ! negCycle[u] {
						negCycle[u] = true
						changed = true
					}
				}
			}
		}
	}

	for i := range(hasPath) {
		if negCycle[i] {
			hasPath[i] = false
		}
	}
	return distances, hasPath, predecessors
}
