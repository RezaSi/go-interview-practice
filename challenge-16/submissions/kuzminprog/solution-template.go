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

// OptimizedSort is your optimized version of SlowSort
// It should produce identical results but perform better
func OptimizedSort(data []int) []int {
	result := make([]int, len(data))
	copy(result, data)

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

// OptimizedStringBuilder is your optimized version of InefficientStringBuilder
// It should produce identical results but perform better
func OptimizedStringBuilder(parts []string, repeatCount int) string {
	var sb strings.Builder
	sb.Grow(len(parts) * repeatCount)

	for i := 0; i < repeatCount; i++ {
		for _, part := range parts {
			sb.WriteString(part)
		}
	}
	return sb.String()
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
	if n <= 1 {
		return n
	}

	first, second := 0, 1
	sum := 0
	for range n {
		sum += second
		first, second = second, first+second
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

	if len(substr) == 0 {
		for i := 0; i < len(text); i++ {
			result[i] = ""
		}
		return result
	}

	subLen := len(substr)

	isASCII := true
	for i := 0; i < subLen; i++ {
		if substr[i] >= 128 {
			isASCII = false
			break
		}
	}

	if isASCII {
		for i := 0; i+subLen <= len(text); i++ {
			match := true

			for j := 0; j < subLen; j++ {
				a := text[i+j]
				b := substr[j]

				if a != b {
					if 'A' <= a && a <= 'Z' {
						a += 'a' - 'A'
					}
					if 'A' <= b && b <= 'Z' {
						b += 'a' - 'A'
					}
					if a != b {
						match = false
						break
					}
				}
			}

			if match {
				result[i] = text[i : i+subLen]
			}
		}
		return result
	}

	lowerSub := strings.ToLower(substr)
	for i := 0; i+subLen <= len(text); i++ {
		if strings.EqualFold(text[i:i+subLen], lowerSub) {
			result[i] = text[i : i+subLen]
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
