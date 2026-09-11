package main

import (
	"slices"
	"strings"
	"sync"
	"time"
)

// SlowSort sorts a slice of integers using a very inefficient algorithm (bubble sort)
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
	result := make([]int, len(data))
	copy(result, data)
	slices.Sort(result)
	return result
}

// InefficientStringBuilder builds a string by repeatedly concatenating
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
	var builder strings.Builder
	for i := 0; i < repeatCount; i++ {
		for _, p := range parts {
			builder.WriteString(p)
		}
	}
	return builder.String()
}

// ExpensiveCalculation performs a computation with redundant work
// It computes the sum of all fibonacci numbers up to n
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

var (
	sumMap sync.Map
	sumMu  sync.Mutex
)

// OptimizedCalculation is your optimized version of ExpensiveCalculation
// It should produce identical results but perform better
func OptimizedCalculation(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return n
	}
	if v, ok := sumMap.Load(n); ok {
		return v.(int)
	}
	result := OptimizedFib(n) + OptimizedCalculation(n-1)
	sumMap.Store(n, result)
	return result
}

var (
	fibMap sync.Map
	fibMu  sync.Mutex
)

func OptimizedFib(n int) int {
	if n <= 1 {
		return n
	} else if v, ok := fibMap.Load(n); ok {
		return v.(int)
	} else {
		result := OptimizedFib(n-1) + OptimizedFib(n-2)
		fibMap.Store(n, result)
		return result
	}
}

// HighAllocationSearch searches for all occurrences of a substring and creates a map with their positions
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
				// Store the 	original case version
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
	if len(substr) == 0 {
		for i := 0; i < len(text); i++ {
			result[i] = ""
		}
		return result
	}

	subLen := len(substr)
	firstCharLower := strings.ToLower(string(substr[0]))
	firstCharUpper := strings.ToUpper(string(substr[0]))

	for i := 0; i <= len(text)-subLen; i++ {
		if text[i:i+1] == firstCharLower || text[i:i+1] == firstCharUpper {
			if strings.EqualFold(text[i:i+subLen], substr) {
				result[i] = text[i : i+subLen]
			}
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
