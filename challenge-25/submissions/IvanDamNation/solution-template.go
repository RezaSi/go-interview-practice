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

const inf = 1000000000

// BreadthFirstSearch implements BFS for unweighted graphs to find shortest paths
// from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func BreadthFirstSearch(graph [][]int, source int) ([]int, []int) {
	dists := make([]int, len(graph))
	parents := make([]int, len(graph))
	for i := range parents {
	    dists[i] = inf
	    parents[i] = -1
	}
	
	queue := make([]int, 0, len(graph))
	queue = append(queue, source)
	dists[source] = 0
	
	ptr := 0
	for ptr < len(queue) {
	    node := queue[ptr]
	    ptr++
	    
	    for _, v := range graph[node] {
	        if dists[v] != inf {
	            continue
	        }
	        
	        parents[v] = node
	        dists[v] = dists[node] + 1
	        queue = append(queue, v)
	    }
	}
	
	return dists, parents
}

// Dijkstra implements Dijkstra's algorithm for weighted graphs with non-negative weights
// to find shortest paths from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func Dijkstra(graph [][]int, weights [][]int, source int) ([]int, []int) {
	n := len(graph)
	dists := make([]int, n)
	parents := make([]int, n)
	for i := range dists {
	    dists[i] = inf
	    parents[i] = -1
	}
	
	visited := make([]bool, n)
	
	h := &MinHeap{
	    nodes: make([]int, 0, n),
	    dists: dists,
	}
	heap.Init(h)
	heap.Push(h, source)
	dists[source] = 0
	
	for h.Len() > 0 {
	    node := heap.Pop(h).(int)
	    
	    if visited[node] {
	        continue
	    }
	    visited[node] = true
	    
	    for i, v := range graph[node] {
	        weight := weights[node][i]
	        
	        if dists[node] + weight < dists[v] {
	            dists[v] = dists[node] + weight
	            parents[v] = node
	            heap.Push(h, v)
	        }
	    }
	}
	
	return dists, parents
}

type MinHeap struct {
    nodes []int
    dists []int
}

func (h MinHeap) Len() int              { return len(h.nodes) }
func (h MinHeap) Less(i, j int) bool    { return h.dists[h.nodes[i]] < h.dists[h.nodes[j]] }
func (h MinHeap) Swap(i, j int)         { h.nodes[i], h.nodes[j] = h.nodes[j], h.nodes[i] }

func (h *MinHeap) Push(item any) {
    h.nodes = append(h.nodes, item.(int))
}
func (h *MinHeap) Pop() any {
    old := h.nodes
    item := old[len(old)-1]
    h.nodes = old[:len(old)-1]
    return item
}

// BellmanFord implements the Bellman-Ford algorithm for weighted graphs that may contain
// negative weight edges to find shortest paths from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - hasPath: slice where hasPath[i] is true if there is a path from source to i without a negative cycle
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func BellmanFord(graph [][]int, weights [][]int, source int) ([]int, []bool, []int) {
	n := len(graph)
	dists := make([]int, n)
	hasPath := make([]bool, n)
	parents := make([]int, n)
	
	for i := range n {
	    dists[i] = inf
	    hasPath[i] = true
	    parents[i] = -1
	}
	dists[source] = 0
	
	for i := 0; i < n-1; i++ {
	    for u := 0; u < n; u++ {
	        if dists[u] == inf {
	            continue
	        }
	        for idx, v := range graph[u] {
	            weight := weights[u][idx]
	            if dists[u] + weight < dists[v] {
	                dists[v] = dists[u] + weight
	                parents[v] = u
	            }
	        }
	    }
	}
	
	for i := 0; i < n; i++ {
	    for u := 0; u < n; u++ {
	        if dists[u] == inf {
	            continue
	        }
	        for idx, v := range graph[u] {
	            weight := weights[u][idx]
	            if dists[u] + weight < dists[v] || !hasPath[u] {
	                if hasPath[v] {
	                    hasPath[v] = false
	                }
	            }
	        }
	    }
	}
	
	for i := range n {
	    if dists[i] == inf {
	        hasPath[i] = false
	    }
	}
	
	return dists, hasPath, parents
}
