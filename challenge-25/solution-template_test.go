package main

import (
	"reflect"
	"testing"
)

func TestBreadthFirstSearch(t *testing.T) {
	tests := []struct {
		name     string
		graph    [][]int
		source   int
		wantDist []int
		wantPred []int
	}{
		{
			name: "Simple unweighted graph",
			graph: [][]int{
				{1, 2},    // Vertex 0 has edges to vertices 1 and 2
				{0, 3, 4}, // Vertex 1 has edges to vertices 0, 3, and 4
				{0, 5},    // Vertex 2 has edges to vertices 0 and 5
				{1},       // Vertex 3 has an edge to vertex 1
				{1},       // Vertex 4 has an edge to vertex 1
				{2},       // Vertex 5 has an edge to vertex 2
			},
			source:   0,
			wantDist: []int{0, 1, 1, 2, 2, 2},
			wantPred: []int{-1, 0, 0, 1, 1, 2},
		},
		{
			name: "Disconnected graph",
			graph: [][]int{
				{1}, // Vertex 0 has an edge to vertex 1
				{0}, // Vertex 1 has an edge to vertex 0
				{3}, // Vertex 2 has an edge to vertex 3
				{2}, // Vertex 3 has an edge to vertex 2
			},
			source:   0,
			wantDist: []int{0, 1, 1000000000, 1000000000},
			wantPred: []int{-1, 0, -1, -1},
		},
		{
			name: "Linear graph",
			graph: [][]int{
				{1}, // Vertex 0 has an edge to vertex 1
				{2}, // Vertex 1 has an edge to vertex 2
				{3}, // Vertex 2 has an edge to vertex 3
				{4}, // Vertex 3 has an edge to vertex 4
				{},  // Vertex 4 has no outgoing edges
			},
			source:   0,
			wantDist: []int{0, 1, 2, 3, 4},
			wantPred: []int{-1, 0, 1, 2, 3},
		},
		{
			name: "Isolated vertex",
			graph: [][]int{
				{1, 2}, // Vertex 0 has edges to vertices 1 and 2
				{0},    // Vertex 1 has an edge to vertex 0
				{0},    // Vertex 2 has an edge to vertex 0
				{},     // Vertex 3 is isolated
			},
			source:   0,
			wantDist: []int{0, 1, 1, 1000000000},
			wantPred: []int{-1, 0, 0, -1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDist, gotPred := BreadthFirstSearch(tt.graph, tt.source)
			if !reflect.DeepEqual(gotDist, tt.wantDist) {
				t.Errorf("BreadthFirstSearch() distances = %v, want %v", gotDist, tt.wantDist)
			}
			if !reflect.DeepEqual(gotPred, tt.wantPred) {
				t.Errorf("BreadthFirstSearch() predecessors = %v, want %v", gotPred, tt.wantPred)
			}
		})
	}
}

func TestDijkstra(t *testing.T) {
	tests := []struct {
		name     string
		graph    [][]int
		weights  [][]int
		source   int
		wantDist []int
		wantPred []int
	}{
		{
			name: "Simple weighted graph",
			graph: [][]int{
				{1, 2},    // Vertex 0 has edges to vertices 1 and 2
				{0, 3, 4}, // Vertex 1 has edges to vertices 0, 3, and 4
				{0, 5},    // Vertex 2 has edges to vertices 0 and 5
				{1},       // Vertex 3 has an edge to vertex 1
				{1},       // Vertex 4 has an edge to vertex 1
				{2},       // Vertex 5 has an edge to vertex 2
			},
			weights: [][]int{
				{5, 10},   // Edge from 0 to 1 has weight 5, edge from 0 to 2 has weight 10
				{5, 3, 2}, // Edge weights from vertex 1
				{10, 2},   // Edge weights from vertex 2
				{3},       // Edge weights from vertex 3
				{2},       // Edge weights from vertex 4
				{2},       // Edge weights from vertex 5
			},
			source:   0,
			wantDist: []int{0, 5, 10, 8, 7, 12},
			wantPred: []int{-1, 0, 0, 1, 1, 2},
		},
		{
			name: "Disconnected weighted graph",
			graph: [][]int{
				{1}, // Vertex 0 has an edge to vertex 1
				{0}, // Vertex 1 has an edge to vertex 0
				{3}, // Vertex 2 has an edge to vertex 3
				{2}, // Vertex 3 has an edge to vertex 2
			},
			weights: [][]int{
				{5}, // Edge from 0 to 1 has weight 5
				{5}, // Edge from 1 to 0 has weight 5
				{2}, // Edge from 2 to 3 has weight 2
				{2}, // Edge from 3 to 2 has weight 2
			},
			source:   0,
			wantDist: []int{0, 5, 1000000000, 1000000000},
			wantPred: []int{-1, 0, -1, -1},
		},
		{
			name: "Linear weighted graph",
			graph: [][]int{
				{1}, // Vertex 0 has an edge to vertex 1
				{2}, // Vertex 1 has an edge to vertex 2
				{3}, // Vertex 2 has an edge to vertex 3
				{4}, // Vertex 3 has an edge to vertex 4
				{},  // Vertex 4 has no outgoing edges
			},
			weights: [][]int{
				{10}, // Edge from 0 to 1 has weight 10
				{20}, // Edge from 1 to 2 has weight 20
				{30}, // Edge from 2 to 3 has weight 30
				{40}, // Edge from 3 to 4 has weight 40
				{},   // No outgoing edges from vertex 4
			},
			source:   0,
			wantDist: []int{0, 10, 30, 60, 100},
			wantPred: []int{-1, 0, 1, 2, 3},
		},
		{
			name: "Zero weight edges",
			graph: [][]int{
				{1, 2}, // Vertex 0 has edges to vertices 1 and 2
				{3},    // Vertex 1 has an edge to vertex 3
				{3},    // Vertex 2 has an edge to vertex 3
				{},     // Vertex 3 has no outgoing edges
			},
			weights: [][]int{
				{0, 0}, // Edges from 0 have zero weight
				{0},    // Edge from 1 has zero weight
				{0},    // Edge from 2 has zero weight
				{},     // No outgoing edges from vertex 3
			},
			source:   0,
			wantDist: []int{0, 0, 0, 0},
			wantPred: []int{-1, 0, 0, 1}, // Path taken is arbitrary when multiple zero-weight paths exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDist, gotPred := Dijkstra(tt.graph, tt.weights, tt.source)
			if !reflect.DeepEqual(gotDist, tt.wantDist) {
				t.Errorf("Dijkstra() distances = %v, want %v", gotDist, tt.wantDist)
			}

			// For zero weight test case, multiple paths could be valid
			if tt.name == "Zero weight edges" {
				// Just check if all distances are correct
				for i, dist := range gotDist {
					if dist != tt.wantDist[i] {
						t.Errorf("Dijkstra() distance[%d] = %v, want %v", i, dist, tt.wantDist[i])
					}
				}
			} else if !reflect.DeepEqual(gotPred, tt.wantPred) {
				t.Errorf("Dijkstra() predecessors = %v, want %v", gotPred, tt.wantPred)
			}
		})
	}
}

func TestBellmanFord(t *testing.T) {
	tests := []struct {
		name        string
		graph       [][]int
		weights     [][]int
		source      int
		wantDist    []int
		wantHasPath []bool
		wantPred    []int
	}{
		{
			name: "Simple graph with negative weights",
			graph: [][]int{
				{1, 2},
				{3},
				{1, 3},
				{4},
				{},
			},
			weights: [][]int{
				{6, 7},  // Edge weights from vertex 0
				{5},     // Edge weights from vertex 1
				{-2, 4}, // Edge weights from vertex 2 (note the negative weight)
				{2},     // Edge weights from vertex 3
				{},      // Edge weights from vertex 4
			},
			source:      0,
			wantDist:    []int{0, 5, 7, 10, 12},
			wantHasPath: []bool{true, true, true, true, true},
			wantPred:    []int{-1, 2, 0, 1, 3},
		},
		{
			name: "Graph with negative cycle",
			graph: [][]int{
				{1}, // Vertex 0 has edge to vertex 1
				{2}, // Vertex 1 has edge to vertex 2
				{3}, // Vertex 2 has edge to vertex 3
				{1}, // Vertex 3 has edge to vertex 1 - creating a cycle 1->2->3->1
			},
			weights: [][]int{
				{2},  // Edge from 0 to 1 has weight 2
				{2},  // Edge from 1 to 2 has weight 2
				{2},  // Edge from 2 to 3 has weight 2
				{-7}, // Edge from 3 to 1 has weight -7 (negative cycle)
			},
			source:      0,
			wantDist:    []int{0, 2, 4, 6},
			wantHasPath: []bool{true, false, false, false}, // Only source vertex has a valid path, others are in negative cycle
			wantPred:    []int{-1, 0, 1, 2},
		},
		{
			name: "Disconnected graph with negative weight",
			graph: [][]int{
				{1}, // Vertex 0 has edge to vertex 1
				{0}, // Vertex 1 has edge to vertex 0
				{3}, // Vertex 2 has edge to vertex 3
				{2}, // Vertex 3 has edge to vertex 2
			},
			weights: [][]int{
				{5},  // Edge from 0 to 1 has weight 5
				{-2}, // Edge from 1 to 0 has weight -2
				{2},  // Edge from 2 to 3 has weight 2
				{-1}, // Edge from 3 to 2 has weight -1
			},
			source:      0,
			wantDist:    []int{0, 5, 1000000000, 1000000000},
			wantHasPath: []bool{true, true, false, false},
			wantPred:    []int{-1, 0, -1, -1},
		},
		{
			name: "Linear graph with mixed weights",
			graph: [][]int{
				{1}, // Vertex 0 has edge to vertex 1
				{2}, // Vertex 1 has edge to vertex 2
				{3}, // Vertex 2 has edge to vertex 3
				{},  // Vertex 3 has no outgoing edges
			},
			weights: [][]int{
				{-5}, // Edge from 0 to 1 has weight -5
				{10}, // Edge from 1 to 2 has weight 10
				{-3}, // Edge from 2 to 3 has weight -3
				{},   // No edges from vertex 3
			},
			source:      0,
			wantDist:    []int{0, -5, 5, 2},
			wantHasPath: []bool{true, true, true, true},
			wantPred:    []int{-1, 0, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDist, gotHasPath, gotPred := BellmanFord(tt.graph, tt.weights, tt.source)

			// For graphs with negative cycles, the exact distances may vary
			// We'll just check if the hasPath boolean is correct for negative cycle cases
			if tt.name == "Graph with negative cycle" {
				if !reflect.DeepEqual(gotHasPath, tt.wantHasPath) {
					t.Errorf("BellmanFord() hasPath = %v, want %v", gotHasPath, tt.wantHasPath)
				}
			} else {
				if !reflect.DeepEqual(gotDist, tt.wantDist) {
					t.Errorf("BellmanFord() distances = %v, want %v", gotDist, tt.wantDist)
				}
				if !reflect.DeepEqual(gotHasPath, tt.wantHasPath) {
					t.Errorf("BellmanFord() hasPath = %v, want %v", gotHasPath, tt.wantHasPath)
				}
				if !reflect.DeepEqual(gotPred, tt.wantPred) {
					t.Errorf("BellmanFord() predecessors = %v, want %v", gotPred, tt.wantPred)
				}
			}
		})
	}
}
