package main

import (
	"fmt"
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

// BreadthFirstSearch implements BFS for unweighted graphs to find shortest paths
// from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func BreadthFirstSearch(graph [][]int, source int) ([]int, []int) {
	distances := make([]int, len(graph))
	predecessors := make([]int, len(graph))
	for i := range graph {
		distances[i] = 1000000000
		predecessors[i] = -1
	}
	queue := make([]int, 0)
	visited := make(map[int]bool)
	queue = append(queue, source)
	visited[source] = true
	distances[source] = 0
	predecessors[source] = -1
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, neighbor := range graph[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
				distances[neighbor] = distances[current] + 1
				predecessors[neighbor] = current
			}
		}
	}

	return distances, predecessors
}

// Dijkstra implements Dijkstra's algorithm for weighted graphs with non-negative weights
// to find shortest paths from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func Dijkstra(graph [][]int, weights [][]int, source int) ([]int, []int) {
	distances := make([]int, len(graph))
	predecessors := make([]int, len(graph))
	for i := range graph {
		distances[i] = 1000000000
		predecessors[i] = -1
	}
	distances[source] = 0
	predecessors[source] = -1
	queue := make([]int, 0)
	queue = append(queue, source)
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for i, neighbor := range graph[current] {
			currentDistance := distances[current] + weights[current][i]
			if currentDistance < distances[neighbor] {
				queue = append(queue, neighbor)
				distances[neighbor] = currentDistance
				predecessors[neighbor] = current
			}
		}
	}
	return distances, predecessors
}

// BellmanFord implements the Bellman-Ford algorithm for weighted graphs that may contain
// negative weight edges to find shortest paths from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - hasPath: slice where hasPath[i] is true if there is a path from source to i without a negative cycle
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func BellmanFord(graph [][]int, weights [][]int, source int) ([]int, []bool, []int) {
	n := len(graph)
	distances := make([]int, n)
	hasPath := make([]bool, n)
	predecessors := make([]int, n)

	for i := 0; i < n; i++ {
		distances[i] = 1000000000
		predecessors[i] = -1
		hasPath[i] = false
	}
	distances[source] = 0
	hasPath[source] = true

	// Relax edges V-1 times
	for i := 0; i < n-1; i++ {
		for u := 0; u < n; u++ {
			if distances[u] != 1000000000 {
				for j, v := range graph[u] {
					if distances[u]+weights[u][j] < distances[v] {
						distances[v] = distances[u] + weights[u][j]
						predecessors[v] = u
						hasPath[v] = true
					}
				}
			}
		}
	}

	// Detect and propagate negative cycles
	queue := make([]int, 0)
	for u := 0; u < n; u++ {
		if distances[u] != 1000000000 {
			for j, v := range graph[u] {
				if distances[u]+weights[u][j] < distances[v] {
					if hasPath[v] {
						hasPath[v] = false
						queue = append(queue, v)
					}
				}
			}
		}
	}

	// Propagate the negative cycle reachable state
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range graph[u] {
			if hasPath[v] {
				hasPath[v] = false
				queue = append(queue, v)
			}
		}
	}

	return distances, hasPath, predecessors
}
