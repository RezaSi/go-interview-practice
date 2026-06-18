package main

import (
	"sort"
	"strings"
	"time"
)

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

// OptimizedSort uses Go's built-in sort package which implements introsort
// (intro sort: a hybrid of quicksort, heapsort, and insertion sort)
// Time complexity: O(n log n) vs O(n²) for bubble sort
func OptimizedSort(data []int) []int {
	// Make a copy to avoid modifying the original
	result := make([]int, len(data))
	copy(result, data)
	
	// Use Go's highly optimized sort function
	sort.Ints(result)
	
	return result
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

// OptimizedStringBuilder uses strings.Builder for efficient string concatenation
// strings.Builder minimizes memory copying by growing its internal buffer
func OptimizedStringBuilder(parts []string, repeatCount int) string {
	// Calculate total length to minimize reallocations
	totalLen := 0
	for _, part := range parts {
		totalLen += len(part)
	}
	totalLen *= repeatCount

	var builder strings.Builder
	builder.Grow(totalLen) // Pre-allocate capacity

	for i := 0; i < repeatCount; i++ {
		for _, part := range parts {
			builder.WriteString(part)
		}
	}

	return builder.String()
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

// OptimizedCalculation uses dynamic programming to avoid redundant calculations
// Instead of recalculating fibonacci numbers, we compute them iteratively
func OptimizedCalculation(n int) int {
	if n <= 0 {
		return 0
	}

	sum := 0
	a, b := 0, 1

	// Compute fibonacci numbers iteratively and sum them
	for i := 1; i <= n; i++ {
		if i == 1 {
			sum += 1
		} else {
			fib := a + b
			sum += fib
			a, b = b, fib
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

// OptimizedSearch reduces allocations by avoiding temporary string creation
// and using more efficient string comparison methods
func OptimizedSearch(text, substr string) map[int]string {
	if len(substr) == 0 {
		return make(map[int]string)
	}

	result := make(map[int]string)
	textLen := len(text)
	substrLen := len(substr)

	// Convert substr to lowercase once
	lowerSubstr := strings.ToLower(substr)

	for i := 0; i <= textLen-substrLen; i++ {
		// Compare characters directly without creating substring
		match := true
		for j := 0; j < substrLen; j++ {
			// Convert to lowercase during comparison to avoid allocation
			textChar := text[i+j]
			if textChar >= 'A' && textChar <= 'Z' {
				textChar += 32 // Convert to lowercase
			}
			
			if textChar != lowerSubstr[j] {
				match = false
				break
			}
		}

		if match {
			// Only allocate when we have a match
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