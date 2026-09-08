package main

import "fmt"

const infinity = int(1e9)

// BreadthFirstSearch implements BFS for unweighted graphs to find shortest paths
// from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func BreadthFirstSearch(graph [][]int, source int) ([]int, []int) {
	n := len(graph)
	distances := make([]int, n)
	predecessors := make([]int, n)

	for i := range distances {
		distances[i] = infinity
		predecessors[i] = -1
	}

	distances[source] = 0

	queue := []int{source}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, neighbor := range graph[curr] {
			if distances[neighbor] == infinity {
				predecessors[neighbor] = curr
				distances[neighbor] = distances[curr] + 1
				queue = append(queue, neighbor)
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
	n := len(graph)

	distances := make([]int, n)
	predecessors := make([]int, n)
	visited := make([]bool, n)

	for i := 0; i < n; i++ {
		distances[i] = infinity
		predecessors[i] = -1
	}
	distances[source] = 0

	for i := 0; i < n; i++ {
		u := -1
		for v := 0; v < n; v++ {
			if !visited[v] && (u == -1 || distances[v] < distances[u]) {
				u = v
			}
		}

		if u == -1 {
			break
		}

		visited[u] = true

		for k, v := range graph[u] {
			weight := weights[u][k]
			if !visited[v] {
				newDist := distances[u] + weight
				if newDist < distances[v] {
					distances[v] = newDist
					predecessors[v] = u
				}
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

	dist := make([]int, n)
	predecessors := make([]int, n)
	onNegCycle := make([]bool, n)
	for i := range dist {
		dist[i] = infinity
		predecessors[i] = -1
		onNegCycle[i] = false
	}
	dist[source] = 0

	for iter := 0; iter < n-1; iter++ {
		for u, neighbor := range graph {
			for j, w := range weights[u] {
				v := neighbor[j]
				if dist[u] != infinity && dist[u]+w < dist[v] {
					dist[v] = dist[u] + w
					predecessors[v] = u
				}
			}
		}
	}

	for u, neighbor := range graph {
		for j, w := range weights[u] {
			v := neighbor[j]
			if dist[u] != infinity && dist[u]+w < dist[v] {
				onNegCycle[v] = true
			}
		}
	}

	queue := make([]int, 0, n)
	for i, affected := range onNegCycle {
		if affected {
			queue = append(queue, i)
		}
	}

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, v := range graph[u] {
			if !onNegCycle[v] {
				onNegCycle[v] = true
				queue = append(queue, v)
			}
		}
	}

	hasPath := make([]bool, n)
	for i := range hasPath {
		hasPath[i] = dist[i] != infinity && !onNegCycle[i]
	}

	return dist, hasPath, predecessors
}

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
