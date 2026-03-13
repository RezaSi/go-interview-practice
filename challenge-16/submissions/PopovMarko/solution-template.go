package main

import (
	"sort"
	"strings"
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
	// Make a copy to avoid modifying the original
	result := make([]int, len(data))
	copy(result, data)
	// Sort using standard library
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
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
func OptimizedStringBuilder(parts []string, repeatCount int) string {
	// Optimized using strings.Builder and bytes.Buffer
	var builder strings.Builder
	for range repeatCount {
		for _, part := range parts {
			builder.Write([]byte(part))
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

// OptimizedCalculation is your optimized version of ExpensiveCalculation
func OptimizedCalculation(n int) int {
	// Optimized by using simple loop with memoization to avoiding redundant calculations
	if n <= 0 {
		return 0
	}
	sum := 1
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
		sum += b
	}
	return sum
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
				// Store the original case version
				result[i] = text[i : i+len(substr)]
			}
		}
	}
	return result
}

// OptimizedSearch is your optimized version of HighAllocationSearch
func OptimizedSearch(text, substr string) map[int]string {
	// Optimized by one time allocations of the variable outside of the loop.
	result := make(map[int]string)
	if text == "" && substr == "" {
		return result
	}
	// Convert to lowercase for case-insensitive search
	lowerText := strings.ToLower(text)
	lowerSubstr := strings.ToLower(substr)
	for i := 0; i <= len(lowerText)-len(lowerSubstr); i++ {
		if lowerText[i:i+len(lowerSubstr)] == lowerSubstr {
			result[i] = text[i : i+len(lowerSubstr)]
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
