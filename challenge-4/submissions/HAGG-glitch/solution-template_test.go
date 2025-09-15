package main

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func bfsReference(graph map[int][]int, start int) []int {
	// A simple reference BFS for checking correctness (sequential).
	queue := []int{start}
	visited := make(map[int]bool)
	visited[start] = true
	var order []int

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)

		for _, v := range graph[u] {
			if !visited[v] {
				visited[v] = true
				queue = append(queue, v)
			}
		}
	}

	return order
}

func buildSampleGraph() map[int][]int {
	// A sample graph
	// 0 -> 1, 2
	// 1 -> 2, 3
	// 2 -> 3
	// 3 -> 4
	// 4 -> (none)
	// 5 -> 2
	return map[int][]int{
		0: {1, 2},
		1: {2, 3},
		2: {3},
		3: {4},
		4: {},
		5: {2},
	}
}

func buildLargeLinearGraph(size int) map[int][]int {
	// Creates a linear graph: 0->1->2->...->size-1
	graph := make(map[int][]int)
	for i := 0; i < size-1; i++ {
		graph[i] = []int{i + 1}
	}
	graph[size-1] = []int{}
	return graph
}

func buildStarGraph(center, branches int) map[int][]int {
	// Creates a star graph with one central node connected to many leaves
	graph := make(map[int][]int)
	centerConnections := make([]int, branches)
	for i := 0; i < branches; i++ {
		leafNode := center + 1 + i
		centerConnections[i] = leafNode
		graph[leafNode] = []int{} // leaf nodes have no connections
	}
	graph[center] = centerConnections
	return graph
}

func TestBasicFunctionality(t *testing.T) {
	graph := buildSampleGraph()

	testCases := []struct {
		name       string
		graph      map[int][]int
		queries    []int
		numWorkers int
	}{
		{
			name:       "Single query, single worker",
			graph:      graph,
			queries:    []int{0},
			numWorkers: 1,
		},
		{
			name:       "Multiple queries, multiple workers",
			graph:      graph,
			queries:    []int{0, 1, 5},
			numWorkers: 2,
		},
		{
			name:       "No queries",
			graph:      graph,
			queries:    []int{},
			numWorkers: 2,
		},
		{
			name:       "More workers than queries",
			graph:      graph,
			queries:    []int{3},
			numWorkers: 10,
		},
		{
			name:       "Many queries, few workers",
			graph:      graph,
			queries:    []int{0, 1, 2, 3, 4, 5},
			numWorkers: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results := ConcurrentBFSQueries(tc.graph, tc.queries, tc.numWorkers)

			if len(results) != len(tc.queries) {
				t.Errorf("Expected %d results, got %d", len(tc.queries), len(results))
				return
			}

			for _, start := range tc.queries {
				refOrder := bfsReference(tc.graph, start)
				if !reflect.DeepEqual(results[start], refOrder) {
					t.Errorf("For start %d, expected %v, got %v", start, refOrder, results[start])
				}
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("Zero workers", func(t *testing.T) {
		graph := buildSampleGraph()
		queries := []int{0, 1}
		results := ConcurrentBFSQueries(graph, queries, 0)
		if len(results) != 0 {
			t.Errorf("Expected empty results when numWorkers=0, but got %v", results)
		}
	})

	t.Run("Empty graph", func(t *testing.T) {
		graph := map[int][]int{}
		queries := []int{0}
		results := ConcurrentBFSQueries(graph, queries, 1)
		if len(results) != 1 || len(results[0]) != 1 || results[0][0] != 0 {
			t.Errorf("Expected [0] for isolated node, got %v", results[0])
		}
	})

	t.Run("Self-loops", func(t *testing.T) {
		graph := map[int][]int{
			0: {0, 1}, // self-loop
			1: {},
		}
		queries := []int{0}
		results := ConcurrentBFSQueries(graph, queries, 1)
		expected := []int{0, 1}
		if !reflect.DeepEqual(results[0], expected) {
			t.Errorf("Expected %v for self-loop graph, got %v", expected, results[0])
		}
	})

	t.Run("Disconnected components", func(t *testing.T) {
		graph := map[int][]int{
			0: {1},
			1: {},
			2: {3},
			3: {},
			4: {},
		}

		queries := []int{0, 2, 4}
		numWorkers := 3
		results := ConcurrentBFSQueries(graph, queries, numWorkers)

		expectedResults := map[int][]int{
			0: {0, 1},
			2: {2, 3},
			4: {4},
		}

		for start, expected := range expectedResults {
			if !reflect.DeepEqual(results[start], expected) {
				t.Errorf("Start=%d, expected %v, got %v", start, expected, results[start])
			}
		}
	})
}

func TestConcurrencyRequirement(t *testing.T) {
	t.Run("Worker count respected", func(t *testing.T) {
		// Create a graph where BFS takes predictable time
		graph := buildStarGraph(0, 100)
		queries := make([]int, 50) // Many queries
		for i := range queries {
			queries[i] = 0
		}

		// Test with different worker counts
		start1 := time.Now()
		_ = ConcurrentBFSQueries(graph, queries, 1)
		duration1 := time.Since(start1)

		start10 := time.Now()
		_ = ConcurrentBFSQueries(graph, queries, 10)
		duration10 := time.Since(start10)

		// With more workers, it should be significantly faster (allowing for some variance)
		if duration10 >= duration1 {
			t.Logf("Warning: 10 workers (%v) not faster than 1 worker (%v). May indicate lack of concurrency.",
				duration10, duration1)
		}
	})
}

func TestPerformanceRequirement(t *testing.T) {
	t.Run("Large graph performance", func(t *testing.T) {
		// Create a graph that would be slow if processed sequentially
		graph := buildLargeLinearGraph(5000)
		queries := []int{0, 1, 2, 3, 4, 1000, 2000, 3000, 4000}
		numWorkers := 4

		start := time.Now()
		results := ConcurrentBFSQueries(graph, queries, numWorkers)
		duration := time.Since(start)

		// Should complete within reasonable time (1 second is generous)
		if duration > 1*time.Second {
			t.Errorf("Performance test failed: took %v, expected < 1s", duration)
		}

		// Verify correctness
		if len(results) != len(queries) {
			t.Errorf("Expected %d results, got %d", len(queries), len(results))
		}

		for _, query := range queries {
			if len(results[query]) == 0 {
				t.Errorf("Empty result for query %d", query)
			}
		}
	})

	t.Run("Sequential vs Concurrent comparison", func(t *testing.T) {
		graph := buildLargeLinearGraph(2000)
		queries := []int{0, 500, 1000, 1500}

		// Measure sequential time (using 1 worker)
		start := time.Now()
		_ = ConcurrentBFSQueries(graph, queries, 1)
		sequentialTime := time.Since(start)

		// Measure concurrent time (using multiple workers)
		start = time.Now()
		_ = ConcurrentBFSQueries(graph, queries, 4)
		concurrentTime := time.Since(start)

		t.Logf("Sequential time: %v, Concurrent time: %v", sequentialTime, concurrentTime)

		// Concurrent should be faster or at least not significantly slower
		// Allow 20% margin for overhead
		if concurrentTime > sequentialTime*12/10 {
			t.Logf("Warning: Concurrent version (%v) not faster than sequential (%v)",
				concurrentTime, sequentialTime)
		}
	})
}

func TestStressTest(t *testing.T) {
	t.Run("Many small queries", func(t *testing.T) {
		graph := buildSampleGraph()
		queries := make([]int, 1000)
		for i := range queries {
			queries[i] = i % 6 // cycle through nodes 0-5
		}

		start := time.Now()
		results := ConcurrentBFSQueries(graph, queries, 10)
		duration := time.Since(start)

		if duration > 500*time.Millisecond {
			t.Errorf("Stress test too slow: %v", duration)
		}

		// Should have results for each unique query (0-5), not total query count
		expectedUniqueResults := 6
		if len(results) != expectedUniqueResults {
			t.Errorf("Expected %d unique results, got %d", expectedUniqueResults, len(results))
		}
	})

	t.Run("Race condition detection", func(t *testing.T) {
		// Run the same test multiple times to catch race conditions
		graph := buildSampleGraph()
		queries := []int{0, 1, 2, 3, 4, 5}
		numWorkers := 5

		var wg sync.WaitGroup
		numIterations := 20
		results := make([]map[int][]int, numIterations)

		for i := 0; i < numIterations; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				results[idx] = ConcurrentBFSQueries(graph, queries, numWorkers)
			}(i)
		}

		wg.Wait()

		// All results should be identical
		for i := 1; i < numIterations; i++ {
			if !reflect.DeepEqual(results[0], results[i]) {
				t.Errorf("Race condition detected: results differ between runs")
				break
			}
		}
	})
}

func TestCorrectnessWithComplexGraphs(t *testing.T) {
	t.Run("Dense graph", func(t *testing.T) {
		// Create a dense graph where every node connects to every other node
		graph := make(map[int][]int)
		numNodes := 10
		for i := 0; i < numNodes; i++ {
			connections := make([]int, 0, numNodes-1)
			for j := 0; j < numNodes; j++ {
				if i != j {
					connections = append(connections, j)
				}
			}
			graph[i] = connections
		}

		queries := []int{0, 5, 9}
		results := ConcurrentBFSQueries(graph, queries, 3)

		// Verify all results contain all nodes (since graph is fully connected)
		for _, query := range queries {
			if len(results[query]) != numNodes {
				t.Errorf("Query %d: expected %d nodes, got %d", query, numNodes, len(results[query]))
				continue
			}
			// First node should be the query node
			if len(results[query]) > 0 && results[query][0] != query {
				t.Errorf("Query %d: first node should be %d, got %d", query, query, results[query][0])
			}
		}
	})

	t.Run("Tree structure", func(t *testing.T) {
		// Binary tree: 0 -> 1,2; 1 -> 3,4; 2 -> 5,6
		graph := map[int][]int{
			0: {1, 2},
			1: {3, 4},
			2: {5, 6},
			3: {},
			4: {},
			5: {},
			6: {},
		}

		queries := []int{0, 1, 2}
		results := ConcurrentBFSQueries(graph, queries, 2)

		// Verify BFS order for tree traversal
		expectedResults := map[int][]int{
			0: {0, 1, 2, 3, 4, 5, 6}, // Level-order traversal
			1: {1, 3, 4},
			2: {2, 5, 6},
		}

		for start, expected := range expectedResults {
			if !reflect.DeepEqual(results[start], expected) {
				t.Errorf("Tree BFS from %d: expected %v, got %v", start, expected, results[start])
			}
		}
	})
}