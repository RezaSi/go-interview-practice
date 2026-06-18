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
	used := make([]bool, len(graph))
	distances := make([]int, len(graph))
	predecessors := make([]int, len(graph))
	
	for i := range graph {
	    distances[i] = int(1e9)
	    predecessors[i] = -1
	}
	
	distances[source] = 0
	
	var queue []int
	queue = append(queue, source)
	for len(queue) > 0 {
	    v := queue[0]
	    queue = queue[1:]
	    if used[v] {
	        continue
	    }
	    used[v] = true;
	    for _, to := range graph[v] {
	        if used[to] {
	            continue
	        }
	        distances[to] = distances[v] + 1
	        predecessors[to] = v
	        queue = append(queue, to)
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
	used := make([]bool, len(graph))
	distances := make([]int, len(graph))
	predecessors := make([]int, len(graph))
	
	for i := range graph {
	    distances[i] = int(1e9)
	    predecessors[i] = -1
	}
	
	distances[source] = 0
	
	for range graph {
	    v := -1
	    for j := range graph {
	        if !used[j] && (v == -1 || distances[j] < distances[v]) {
	            v = j
	        }
	    }
	    if distances[v] == int(1e9) {
	        break
	    }
	    used[v] = true
	    
	    for j, to := range graph[v] {
	        if distances[to] > distances[v] + weights[v][j] {
	            distances[to] = distances[v] + weights[v][j]
	            predecessors[to] = v
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
    type Edge struct {
        a, b, c int
    }

	var edges []Edge
	for i := range graph {
	    for j, to := range graph[i] {
	        edges = append(edges, Edge{
	            a: i,
	            b: to,
	            c: weights[i][j],
	        })
	    }
	}
	
	used := make([]bool, len(graph))
	distances := make([]int, len(graph))
	predecessors := make([]int, len(graph))
	
	for i := range graph {
	    distances[i] = int(1e9)
	    predecessors[i] = -1
	}
	
	distances[source] = 0
	used[source] = true
	
	var x int
	for range graph {
	    x = -1
	    for _, e := range edges {
	        if distances[e.a] < int(1e9) {
	            if distances[e.b] > distances[e.a] + e.c {
	                distances[e.b] = max(-int(1e9), distances[e.a] + e.c)
	                predecessors[e.b] = e.a
	                used[e.b] = true
	                x = e.b
	            }
	        }
	    }
	}
	
	if x != -1 {
	    y := x
	    for range graph {
	        y = predecessors[y]
	    }
	    curr := y
	    var negative_cycle []int
	    for {
	        negative_cycle = append(negative_cycle, curr)
	        curr = predecessors[curr]
	        if curr == y && len(negative_cycle) > 1 {
	            break
	        }
	    }
	    for _, v := range negative_cycle {
	        used[v] = false
	    }
	}
	return distances, used, predecessors
}
