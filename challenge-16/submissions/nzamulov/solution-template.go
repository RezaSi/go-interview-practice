package main

import (
	"strings"
	"sort"
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
	var fullParts strings.Builder
	for _, part := range parts {
		fullParts.WriteString(part)
	}

    var sb strings.Builder
	var str = fullParts.String()
	sb.Grow(len(str) * repeatCount)
	for i := 0; i < repeatCount; i++ {
		sb.WriteString(str)
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
    if n <= 0 {
        return 0
    }
    if n == 1 {
        return 1
    }
    fibPrev, fibNext := 0, 1
    sum := 1
    for i := 2; i <= n; i++ {
        next := fibPrev + fibNext
        fibPrev, fibNext = fibNext, next
        sum += next
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
    if len(text) == 0 || len(substr) == 0 {
        return result
    }
	lowerText := strings.ToLower(text)
	lowerSubstr := strings.ToLower(substr)
    pattern := lowerSubstr + "#" + lowerText
	p := make([]int, len(pattern))
	for i := 1; i < len(pattern); i++ {
	    j := p[i - 1]
	    for j > 0 && pattern[i] != pattern[j] {
	        j = p[j - 1]
	    }
	    if pattern[i] == pattern[j] {
	        j++
	    }
	    p[i] = j
	}
	for i := range p {
	    if p[i] == len(lowerSubstr) {
	        index := i - 2 * len(lowerSubstr)
	        result[index] = text[index : index + len(lowerSubstr)]
	    }
	}
	return result
}