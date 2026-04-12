package main

import (
	"fmt"
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

// OptimizedSort is your optimized version of SlowSort
// It should produce identical results but perform better
func OptimizedSort(data []int) []int {
	// TODO: Implement a more efficient sorting algorithm
	// Hint: Consider using sort package or a more efficient algorithm
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
	// TODO: Implement a more efficient string building method
	// Hint: Consider using strings.Builder or bytes.Buffer
	var builder strings.Builder
	totalLen := 0
	for _, part := range parts {
		totalLen += len(part)
	}
	builder.Grow(totalLen * repeatCount)

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

// OptimizedCalculation is your optimized version of ExpensiveCalculation
// It should produce identical results but perform better
func OptimizedCalculation(n int) int {
	// TODO: Implement a more efficient calculation method
	// Hint: Consider memoization or avoiding redundant calculations

	if n <= 0 {
		return 0
	}

	sum := 0
	for i := 1; i <= n; i++ {
		sum += optimizedFibonacci(i)
	}

	return sum

}

func optimizedFibonacci(n int) int {
	if n <= 1 {
		return n
	}

	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
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
	// TODO: Implement a more efficient search method with fewer allocations
	// Hint: Consider avoiding temporary string allocations and reusing memory
	result := make(map[int]string)

	// Convert to lowercase for case-insensitive search
	lowerText := strings.ToLower(text)
	lowerSubstr := strings.ToLower(substr)
	lowerTextLen := len(lowerText)
	lowerSubstrLen := len(lowerSubstr)

	if lowerTextLen == 0 || lowerSubstrLen == 0 {
		return result
	}

	for i := 0; i <= lowerTextLen-lowerSubstrLen; {
		idx := strings.Index(lowerText[i:], lowerSubstr)
		if idx == -1 {
			break
		}
		realIdx := i + idx
		result[realIdx] = text[realIdx : realIdx+len(substr)]
		i = realIdx + lowerSubstrLen
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

func main() {
	fmt.Println(HighAllocationSearch("Hello World Hello", "hello"))
	fmt.Println(OptimizedSearch("Hello World Hello", "hello"))
	fmt.Println(HighAllocationSearch("", ""))
	fmt.Println(OptimizedSearch("", ""))
}
