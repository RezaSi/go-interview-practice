package main

import (
	"bytes"
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

// partition partially sorts the array.
// items smaller than pivot are moved to left side of pivot.
// items greater than pivot are moved to riht side of pivot.
func partition(data []int, i, j int) int {
	pivot := data[i]
	x := i
	for p := i + 1; p <= j; p++ {
		if data[p] < pivot {
			x++
			data[p], data[x] = data[x], data[p]
		}
	}

	data[i], data[x] = data[x], data[i]
	return x
}

// quickSort sorts a slice of integer partition smaller items on left side and
// bigger items on right side of pivot element selected in each iteration till slice is sorted.
func quickSort(data []int, start, end int) {
	if start >= end {
		return
	}

	stack := make([]int, end-start+1)
	top := -1

	top++
	stack[top] = start
	top++
	stack[top] = end

	for top >= 0 {
		j := stack[top]
		top--
		i := stack[top]
		top--

		pivot := partition(data, i, j)

		// If element are present on left side of pivot
		if pivot-1 > i {
			top++
			stack[top] = i
			top++
			stack[top] = pivot - 1
		}

		//  if element are present on right side of pivot
		if pivot+1 < j {
			top++
			stack[top] = pivot + 1
			top++
			stack[top] = j
		}
	}
}

// OptimizedSort is your optimized version of SlowSort
// It should produce identical results but perform better
func OptimizedSort(data []int) []int {
	result := make([]int, len(data))
	copy(result, data)
	quickSort(result, 0, len(data)-1)
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
	result := bytes.Buffer{}
	defer result.Reset()

	temp := []byte(strings.Join(parts, ""))

	for range repeatCount {
		result.Write(temp)
	}

	return result.String()
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

// myFib computes the fibbonaci num at n
//
// myFib(1) = 1
//
// myFib(0) = 0
func myFib(n int) int {
	if n <= 1 {
		return n
	}
	a := 1
	b := 0
	result := 1

	for i := 2; i <= n; i++ {
		result = a + b
		b = a
		a = result
	}
	return result
}

// OptimizedCalculation is your optimized version of ExpensiveCalculation
// It should produce identical results but perform better
func OptimizedCalculation(n int) int {

	if n <= 0 {
		return 0
	}

	sum := 0
	for i := 1; i <= n; i++ {
		sum += myFib(i)
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
		return result
	}

	// We lowercase the needle once, as it's usually small.
	// We only lowercase the portion of the haystack we are checking
	// to avoid massive heap allocations for large 'text' inputs.
	lowerSubstr := strings.ToLower(substr)
	substrLen := len(substr)
	limit := len(text) - substrLen

	for i := 0; i <= limit; i++ {
		// Use EqualFold for a highly optimized, allocation-free
		// case-insensitive comparison of the current segment.
		if strings.EqualFold(text[i:i+substrLen], lowerSubstr) {
			result[i] = text[i : i+substrLen]
			// Optional: Skip ahead to avoid overlapping matches
			// i += substrLen - 1
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
