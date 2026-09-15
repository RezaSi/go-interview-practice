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
	gmap := make(map[int][]int)
	visited := map[int]bool{}
	for i, g := range graph {
		gmap[i] = g
		visited[i] = false
	}

	distances := make([]int, len(graph), len(graph))
	predcessors := make([]int, len(graph), len(graph))
	dmap := map[int]int{}
	dmap[source] = 0
	pmap := map[int]int{}
	pmap[source] = -1

	queue := []int{source}
	visited[source] = true

	for {
		if len(queue) == 0 {
			break
		}

		last := len(queue) - 1

		q := queue[last] // extract last

		queue = queue[:last]

		for _, vx := range gmap[q] {
			if visited[vx] != true {
				visited[vx] = true
				queue = append(queue, vx)
				dmap[vx] = dmap[q] + 1 // !!!
				pmap[vx] = q
			}
		}
	}

	for k, v := range dmap {
		distances[k] = v
	}

	for k, v := range pmap {
		predcessors[k] = v
	}

	for k, v := range visited {
		if !v {
			distances[k] = int(1e9)
			predcessors[k] = -1
		}
	}

	return distances, predcessors
}

// Dijkstra implements Dijkstra's algorithm for weighted graphs with non-negative weights
// to find shortest paths from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func Dijkstra(graph [][]int, weights [][]int, source int) ([]int, []int) {
	gmap := make(map[int][]int)
	visited := make(map[int]bool)

	for i, g := range graph {
		gmap[i] = g
		visited[i] = false
	}

	wmap := make(map[int][]int)
	for i, w := range weights {
		wmap[i] = w
	}

	distances := make([]int, len(graph), len(graph))
	predcessors := make([]int, len(graph), len(graph))

	for i := 0; i < len(graph); i++ {
		distances[i] = int(1e9)
		predcessors[i] = -1
	}

	distances[source] = 0

	queue := []int{source}
	visited[source] = true

	var q int

	for {
		if len(queue) == 0 {
			break
		}

		queue, q = extractMin(queue, distances) //[], 0

		for i, v := range wmap[q] { // 5 10
			nv := gmap[q][i]
			if visited[nv] != true {
				if distances[q]+v < distances[nv] {
					distances[nv] = distances[q] + v
					queue = append(queue, nv)
					predcessors[nv] = q
				}
				visited[i] = true
			}
		}

	}

	return distances, predcessors
}

func extractMin(q []int, distances []int) ([]int, int) {
	var min, minIdx int
	var first bool = true
	for i, v := range q {
		if first {
			min = distances[v]
			minIdx = i
			first = false
		}
		if min > distances[v] {
			min = distances[v]
			minIdx = i
		}
	}

	minVertex := q[minIdx]

	q = append(q[:minIdx], q[minIdx+1:]...)

	return q, minVertex

}

// BellmanFord implements the Bellman-Ford algorithm for weighted graphs that may contain
// negative weight edges to find shortest paths from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - hasPath: slice where hasPath[i] is true if there is a path from source to i without a negative cycle
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func BellmanFord(graph [][]int, weights [][]int, source int) ([]int, []bool, []int) {
	gmap := make(map[int][]int)
	for i, g := range graph {
		gmap[i] = g
	}
	wmap := make(map[int][]int)
	for i, w := range weights {
		wmap[i] = w
	}

	const INF = int(1e9)

	distances := make([]int, len(graph), len(graph))
	predcessors := make([]int, len(graph), len(graph))
	hasPath := make([]bool, len(graph), len(graph))

	for i := 0; i < len(graph); i++ {
		distances[i] = INF
		hasPath[i] = true
		predcessors[i] = -1
	}

	distances[source] = 0

	for range len(graph) - 1 {
		for k, wei := range wmap {
			for i, w := range wei {
				if distances[k] != INF && distances[k]+w < distances[gmap[k][i]] {
					distances[gmap[k][i]] = distances[k] + w
					predcessors[gmap[k][i]] = k
				}
			}
		}
	}

	nega := make([]bool, len(graph), len(graph))

	for k, wei := range wmap {
		for i, w := range wei {
			if distances[k] != INF && distances[k]+w < distances[gmap[k][i]] {
				nega[gmap[k][i]] = true
			}
		}
	}

	for range len(graph) - 1 {
		for k, wei := range wmap {
			for i, _ := range wei {
				if nega[k] {
					nega[gmap[k][i]] = true
				}
			}
		}
	}

	for i, d := range distances {
		if d == INF || nega[i] == true {
			hasPath[i] = false
		}
	}

	return distances, hasPath, predcessors
}
