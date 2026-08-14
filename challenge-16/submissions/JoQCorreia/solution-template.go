package main

import (
	"strings"
	"time"
	"slices"
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
	
	sort := make([]int, len(data))
	copy(sort, data)
	slices.Sort(sort)
	return sort
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
	var str strings.Builder
	for i := 0; i < repeatCount; i++ {
		for _, word := range parts {
			str.WriteString(word)
		}
	}
	return str.String() 
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

	fib := []int{0, 1, 0}

	for i := 1; i <= n; i++ {
		if i == 1 {
			sum += (fib[0] + fib[1])
			continue
		}
		fib[2] = fib[0] + fib[1]

		sum += fib[2]

		fib[0] = fib[1]
		fib[1] = fib[2]

	}

	return sum // Replace this with your optimized implementation
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

func lps(pat string) []int {
	leng := 0
	lps := make([]int, len(pat))
	lps[0] = 0
	lowerPat := strings.ToLower(pat)

	for i := 1; i < len(lowerPat); {
		if lowerPat[i] == lowerPat[leng] {
			leng++
			lps[i] = leng
			i++
		} else {
			if leng != 0 {
				leng = lps[leng-1]
			} else {
				lps[i] = 0
				i++
			}

		}
	}

	return lps
}

func OptimizedSearch(text, substr string) map[int]string {

	result := make(map[int]string)

	if len(text) == 0 || len(substr) == 0 {
		return result
	}

	lowerText := strings.ToLower(text)
	lowerSubstr := strings.ToLower(substr)

	lps := lps(substr)
	j := 0

	for i := 0; i < len(lowerText); {
		if lowerText[i] == lowerSubstr[j] {
			i++
			j++

			if j == len(lowerSubstr) {
				result[i-j] = text[i-len(substr) : i]
				j = lps[j-1]
			}
		} else {
			if j != 0 {
				j = lps[j-1]
			} else {
				i++
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
