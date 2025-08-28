package main

import (
	"fmt"
)

const inf = 1000000000

type minHeap struct {
	val, index, loc []int
	size            int
}

func (h *minHeap) bubbleDown(ind int) {
	L, R := inf, inf

	if ind*2+1 < h.size {
		L = h.val[ind*2+1]
	}
	if ind*2+2 < h.size {
		R = h.val[ind*2+2]
	}

	if h.val[ind] <= L && h.val[ind] <= R {
		return
	}

	if L < R {
		h.loc[h.index[ind]], h.loc[h.index[2*ind+1]] = 2*ind+1, ind
		h.val[ind], h.val[2*ind+1] = h.val[2*ind+1], h.val[ind]
		h.index[ind], h.index[2*ind+1] = h.index[2*ind+1], h.val[ind]
		h.bubbleDown(2*ind + 1)
	} else {
		h.loc[h.index[ind]], h.loc[h.index[2*ind+2]] = 2*ind+2, ind
		h.val[ind], h.val[2*ind+2] = h.val[2*ind+2], h.val[ind]
		h.index[ind], h.index[2*ind+2] = h.index[2*ind+2], h.val[ind]
		h.bubbleDown(2*ind + 2)
	}
}

func (h *minHeap) bubbleUp(ind int) {
	if ind == 0 {
		return
	}

	parInd := (ind - 1) / 2
	if h.val[ind] < h.val[parInd] {
		h.loc[h.index[ind]], h.loc[h.index[parInd]] = parInd, ind
		h.val[ind], h.val[parInd] = h.val[parInd], h.val[ind]
		h.index[ind], h.index[parInd] = h.index[parInd], h.index[ind]
		h.bubbleUp(parInd)
	}
}

func (h *minHeap) insert(num, ind int) {
	h.loc[ind] = h.size
	if len(h.val) == h.size {
		h.val = append(h.val, num)
		h.index = append(h.index, ind)
	} else {
		h.val[h.size] = num
		h.index[h.size] = ind
	}
	h.bubbleUp(h.size)
	h.size++
}

func (h *minHeap) remove(v int) {
	ind := h.loc[v]
	h.loc[v] = -1
	h.size--
	if ind == h.size {
		return
	}
	h.loc[h.index[h.size]] = ind
	h.val[ind], h.val[h.size] = h.val[h.size], h.val[ind]
	h.index[ind], h.index[h.size] = h.index[h.size], h.index[ind]

	if ind > 0 && h.val[(ind-1)/2] > h.val[ind] {
		h.bubbleUp(ind)
	} else {
		h.bubbleDown(ind)
	}
}

var h = minHeap{
	val:   make([]int, 0),
	index: make([]int, 0),
	size:  0,
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

// BreadthFirstSearch implements BFS for unweighted graphs to find shortest paths
// from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func BreadthFirstSearch(graph [][]int, source int) ([]int, []int) {
	mark := make([]bool, len(graph))
	dis := make([]int, len(graph))
	pre := make([]int, len(graph))
	queue := make([]int, 0)
	index := 0
	var current int

	for i := 0; i < len(graph); i++ {
		pre[i] = -1
		dis[i] = inf
	}
	dis[source] = 0
	mark[source] = true
	queue = append(queue, source)

	for index < len(queue) {
		current = queue[index]
		index++
		for _, v := range graph[current] {
			if !mark[v] {
				mark[v] = true
				pre[v] = current
				dis[v] = dis[current] + 1
				queue = append(queue, v)
			}
		}
	}

	return dis, pre
}

// Dijkstra implements Dijkstra's algorithm for weighted graphs with non-negative weights
// to find shortest paths from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func Dijkstra(graph [][]int, weights [][]int, source int) ([]int, []int) {
	dis := make([]int, len(graph))
	pre := make([]int, len(graph))
	h.loc = make([]int, len(graph))

	for i := 0; i < len(graph); i++ {
		pre[i] = -1
		h.loc[i] = -1
		dis[i] = inf
	}
	dis[source] = 0
	h.insert(0, source)

	var u int
	for h.size > 0 {
		u = h.index[0]
		h.remove(u)
		for i, v := range graph[u] {
			if dis[u]+weights[u][i] < dis[v] {
				if h.loc[v] != -1 {
					h.remove(v)
				}
				dis[v] = dis[u] + weights[u][i]
				pre[v] = u
				h.insert(dis[v], v)
			}
		}
	}

	return dis, pre
}

// BellmanFord implements the Bellman-Ford algorithm for weighted graphs that may contain
// negative weight edges to find shortest paths from a source vertex to all other vertices.
// Returns:
// - distances: slice where distances[i] is the shortest distance from source to vertex i
// - hasPath: slice where hasPath[i] is true if there is a path from source to i without a negative cycle
// - predecessors: slice where predecessors[i] is the vertex that comes before i in the shortest path
func BellmanFord(graph [][]int, weights [][]int, source int) ([]int, []bool, []int) {
	dis := make([]int, len(graph))
	has := make([]bool, len(graph))
	pre := make([]int, len(graph))

	for i := 0; i < len(graph); i++ {
		pre[i] = -1
		dis[i] = 2 * inf
		has[i] = true
	}
	dis[source] = 0

	for i := 0; i < len(graph)-1; i++ {
		for u := 0; u < len(graph); u++ {
			for j, v := range graph[u] {
				if dis[u]+weights[u][j] < dis[v] {
					dis[v] = dis[u] + weights[u][j]
					pre[v] = u
				}
			}
		}
	}

	var currentV int
	for u := 0; u < len(graph); u++ {
		for i, v := range graph[u] {
			if dis[u]+weights[u][i] < dis[v] {
				has[v] = false
				currentV = u
				for currentV != v {
					has[currentV] = false
					currentV = pre[currentV]
				}
				break
			}
		}
	}
	for i := 0; i < len(graph); i++ {
		if dis[i] > inf {
			dis[i] = inf
			has[i] = false
			pre[i] = -1
		}
	}

	return dis, has, pre
}
