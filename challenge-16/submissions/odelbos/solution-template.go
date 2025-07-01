package main

import (
	"strings"
	"time"
	"sort"
)

//
// Benchmark:
//
// % go test -bench=. -benchmem
// goos: darwin
// goarch: amd64
// pkg: challenge
// cpu: Intel(R) Core(TM) i7-6700K CPU @ 4.00GHz
// BenchmarkSlowSort-8                       149242       7630 ns/op           896 B/op       1 allocs/op
// BenchmarkOptimizedSort-8                11053554        160.5 ns/op           0 B/op       0 allocs/op
// BenchmarkInefficientStringBuilder-8       356548       3427 ns/op          5424 B/op      69 allocs/op
// BenchmarkOptimizedStringBuilder-8        3098140        338.7 ns/op         144 B/op       1 allocs/op
// BenchmarkExpensiveCalculation-8          1749090        678.7 ns/op           0 B/op       0 allocs/op
// BenchmarkOptimizedCalculation-8         36544660         30.56 ns/op          0 B/op       0 allocs/op
// BenchmarkHighAllocationSearch-8            36944      33573 ns/op         12934 B/op      10 allocs/op
// BenchmarkOptimizedSearch-8                 46342      25345 ns/op          8072 B/op       9 allocs/op
//

// SlowSort sorts a slice of integers using a very inefficient algorithm (bubble sort)
// TODO: Optimize this function to be more efficient
func SlowSort(data []int) []int {
	// Make a copy to avoid modifying the original
	result := make([]int, len(data))
	copy(result, data)

	// Bubble sort implementation
	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result)-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}

	return result
}

// OptimizedSort is your optimized version of SlowSort
// It should produce identical results but perform better
func OptimizedSort(data []int) []int {
    sort.Ints(data)
    return data
}

// InefficientStringBuilder builds a string by repeatedly concatenating
// TODO: Optimize this function to be more efficient
func InefficientStringBuilder(parts []string, repeatCount int) string {
	result := ""

	for i := 0; i < repeatCount; i++ {
		for _, part := range parts {
			result += part
		}
	}

	return result
}

// OptimizedStringBuilder is your optimized version of InefficientStringBuilder
// It should produce identical results but perform better
func OptimizedStringBuilder(parts []string, repeatCount int) string {
    var buffer strings.Builder
    size := 0
    for _, part := range(parts) {
        size += len(part)
    }
    buffer.Grow(size * repeatCount)

	for i := 0; i < repeatCount; i++ {
        for _, part := range(parts) {
            buffer.WriteString(part)
        }
	}
	return buffer.String()
}

// ExpensiveCalculation performs a computation with redundant work
// It computes the sum of all fibonacci numbers up to n
// TODO: Optimize this function to be more efficient
func ExpensiveCalculation(n int) int {
	if n <= 0 {
		return 0
	}

	sum := 0
	for i := 1; i <= n; i++ {
		sum += fibonacci(i)
	}

	return sum
}

// Helper function that computes the fibonacci number at position n
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// OptimizedCalculation is your optimized version of ExpensiveCalculation
// It should produce identical results but perform better
func OptimizedCalculation(n int) int {
	if n <= 0 {
		return 0
	}
	sum := 0
	for i := 1; i <= n; i++ {
        if i <= 1 {
            sum += i
        } else {
            a, b := 0, 1
	        for j := 2; j <= i; j++ {
	            a, b = b, a + b
	        }
	        sum += b
        }
	}
    return sum
}

// HighAllocationSearch searches for all occurrences of a substring and creates a map with their positions
// TODO: Optimize this function to reduce allocations
func HighAllocationSearch(text, substr string) map[int]string {
	result := make(map[int]string)

	// Convert to lowercase for case-insensitive search
	lowerText := strings.ToLower(text)
	lowerSubstr := strings.ToLower(substr)

	for i := 0; i < len(lowerText); i++ {
		// Check if we can fit the substring starting at position i
		if i+len(lowerSubstr) <= len(lowerText) {
			// Extract the potential match
			potentialMatch := lowerText[i : i+len(lowerSubstr)]

			// Check if it matches
			if potentialMatch == lowerSubstr {
				// Store the original case version
				result[i] = text[i : i+len(substr)]
			}
		}
	}

	return result
}

// OptimizedSearch is your optimized version of HighAllocationSearch
// It should produce identical results but perform better with fewer allocations
func OptimizedSearch(text, substr string) map[int]string {
    result := make(map[int]string)
    substrLen := len(substr)
    if substrLen == 0 || len(text) < substrLen {
        return result
    }
    for i := 0; i <= len(text)-substrLen; i++ {
        if strings.EqualFold(text[i:i+substrLen], substr) {
            result[i] = text[i : i+substrLen]
        }
    }
    return result
}

// A function to simulate CPU-intensive work for benchmarking
// You don't need to optimize this; it's just used for testing
func SimulateCPUWork(duration time.Duration) {
	start := time.Now()
	for time.Since(start) < duration {
		// Just waste CPU cycles
		for i := 0; i < 1000000; i++ {
			_ = i
		}
	}
}
