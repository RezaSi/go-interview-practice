package main

import (
	"container/heap"
	"fmt"
)

func main() {
	// Example 1: Unweighted graph for BFS
	unweightedGraph := [][]int{
		{1, 2},
		{0, 3, 4},
		{0, 5},
		{1},
		{1},
		{2},
	}
	distances, predecessors := BreadthFirstSearch(unweightedGraph, 0)
	fmt.Println("BFS Results:")
	fmt.Printf("Distances: %v\n", distances)
	fmt.Printf("Predecessors: %v\n", predecessors)
	fmt.Println()

	// Example 2: Weighted graph for Dijkstra
	weightedGraph := [][]int{
		{1, 2},
		{0, 3, 4},
		{0, 5},
		{1},
		{1},
		{2},
	}
	weights := [][]int{
		{5, 10},   // Edge from 0 to 1 has weight 5, edge from 0 to 2 has weight 10
		{5, 3, 2}, // Edge weights from vertex 1
		{10, 2},   // Edge weights from vertex 2
		{3},       // Edge weights from vertex 3
		{2},       // Edge weights from vertex 4
		{2},       // Edge weights from vertex 5
	}
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
	bfDistances, hasPath, bfPredecessors := BellmanFord(negativeWeightGraph, negativeWeights, 0)
	fmt.Println("Bellman-Ford Results:")
	fmt.Printf("Distances: %v\n", bfDistances)
	fmt.Printf("Has Path: %v\n", hasPath)
	fmt.Printf("Predecessors: %v\n", bfPredecessors)
}

func BreadthFirstSearch(graph [][]int, source int) ([]int, []int) {
	n := len(graph)
	distances := make([]int, n)
	predecessors := make([]int, n)
	visited := make([]bool, n)
	for i := range distances {
		distances[i] = 1e9
		predecessors[i] = -1
	}
	distances[source] = 0
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

type PriorityQueue []Node
type Node struct {
	vertex   int
	distance int
}

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].distance < pq[j].distance }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
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
	inNegativeCycle := make([]bool, n)
	for u := 0; u < n; u++ {
		if distances[u] == 1e9 {
			continue
		}
		for j, v := range graph[u] {
			weight := weights[u][j]
			if distances[u]+weight < distances[v] {
				inNegativeCycle[v] = true
			}
		}
	}
	// Propagate negative cycle reachability using BFS
	queue := []int{}
	for i := 0; i < n; i++ {
		if inNegativeCycle[i] {
			queue = append(queue, i)
			hasPath[i] = false
		}
	}
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
