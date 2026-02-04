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

const (
	inf = 1000000000
)

// BreadthFirstSearch implements BFS for unweighted graphs to find shortest paths
// from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func BreadthFirstSearch(graph [][]int, source int) ([]int, []int) {
	// TODO: Implement this function

	l := len(graph)
	distances := make([]int, l)
	predecessors := make([]int, l)

	mark := make(map[int]bool)

	for i := 0; i < l; i++ {
		distances[i] = inf
		predecessors[i] = -1
	}

	queue := []int{source}
	mark[source] = true
	distances[source] = 0

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		for _, v := range graph[curr] {
			if _, ok := mark[v]; !ok {
				mark[v] = true
				distances[v] = distances[curr] + 1
				predecessors[v] = curr
				queue = append(queue, v)
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
	// TODO: Implement this function
	l := len(graph)
	distances := make([]int, l)
	predecessors := make([]int, l)

	mark := make(map[int]bool)

	for i := 0; i < l; i++ {
		distances[i] = inf
		predecessors[i] = -1
	}

	queue := []int{source}
	mark[source] = true
	distances[source] = 0

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		for i, v := range graph[curr] {
			if _, ok := mark[v]; !ok {
				mark[v] = true
				distances[v] = distances[curr] + weights[curr][i]
				predecessors[v] = curr
				queue = append(queue, v)
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
	// TODO: Implement this function
	l := len(graph)
	distances := make([]int, l)
	predecessors := make([]int, l)
	hasPath := make([]bool, l)

	for i := 0; i < l; i++ {
		distances[i] = inf
		predecessors[i] = -1
		hasPath[i] = false
	}

	distances[source] = 0
	hasPath[source] = true

	for i := 0; i < l-1; i++ {
		changed := false
		for u := 0; u < l; u++ {
			if distances[u] == inf {
				continue
			}
			for j, v := range graph[u] {
				w := weights[u][j]
				if distances[u]+w < distances[v] {
					distances[v] = distances[u] + w
					predecessors[v] = u
					hasPath[v] = true
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}

	visited := make([]bool, l)
	for u := 0; u < l; u++ {
		if distances[u] == inf {
			continue
		}
		for j, v := range graph[u] {
			w := weights[u][j]
			if distances[u]+w < distances[v] {
				markReachableFromCycle(graph, u, distances, hasPath, visited)
			}
		}
	}

	return distances, hasPath, predecessors
}

func markReachableFromCycle(graph [][]int, start int, distances []int, hasPath []bool, visited []bool) {

	for i := range visited {
		visited[i] = false
	}

	queue := []int{start}
	visited[start] = true

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		distances[u] = -inf
		hasPath[u] = false

		for _, v := range graph[u] {
			if !visited[v] {
				visited[v] = true
				queue = append(queue, v)
			}
		}
	}
}
