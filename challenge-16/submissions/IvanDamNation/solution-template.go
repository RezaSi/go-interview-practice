package main

import (
    "cmp"
    "slices"
    "strconv"
	"strings"
	"testing"
	"time"
	"unicode"
	"unicode/utf8"
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
	cData := make([]int, len(data))
	copy(cData, data)
	if len(cData) <= 1 {
	    return cData
	}
	
	slices.SortFunc(cData, func(i, j int) int {
	    return cmp.Compare(i, j)
	})
	
	return cData
}

// BenchmarkSortingAlgorithms starts benchmark for Sort part of challenge
func BenchmarkSortingAlgorithms(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		inputData := generateRandomSlice(size)

		b.Run("SlowSort-Size-"+strconv.Itoa(size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = SlowSort(inputData)
			}
		})

		b.Run("OptimizedSort-Size-"+strconv.Itoa(size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = OptimizedSort(inputData)
			}
		})
	}
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
	if repeatCount <= 0 || len(parts) == 0 {
	    return ""
	}
	
	totalLen := 0
	for _, part := range parts {
	    totalLen += len(part)
	}
	totalLen *= repeatCount
	
	var sb strings.Builder
	sb.Grow(totalLen)
	
	for range repeatCount {
	    for _, part := range parts {
	        sb.WriteString(part)
	    }
	}
	
	return sb.String()
}

// benchStringResult defends from compile optimizations on BenchmarkStringBuilders
var benchStringResult string

// BenchmarkStringBuilders starts benchmark for StringBuilder part of challenge
func BenchmarkStringBuilders(b *testing.B) {
	parts := []string{"hello", "world", "golang", "benchmark", "performance", "optimization"}
	repeatCount := 100

	b.Run("Inefficient", func(b *testing.B) {
		var r string
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r = InefficientStringBuilder(parts, repeatCount)
		}
		benchStringResult = r
	})

	b.Run("Optimized", func(b *testing.B) {
		var r string
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r = OptimizedStringBuilder(parts, repeatCount)
		}
		benchStringResult = r
	})
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
// It should produce identical results but perform better
func OptimizedCalculation(n int) int {
	if n <= 0 {
	    return 0
	}
	
	if n == 1 {
	    return 1
	}
	
	prev, fib := 1, 1
	for range n {
	    next := prev + fib
	    prev, fib = fib, next
	}
	return fib-1
}

// benchResult defends from compile optimizations on BenchmarkCalculationComparison
var benchResult int

// BenchmarkCalculationComparison starts benchmark for Calculation part of challenge
func BenchmarkCalculationComparison(b *testing.B) {
	var r int

	b.Run("Expensive-n-10", func(b *testing.B) {
		for i := 0; i < b.N; i++ { r = ExpensiveCalculation(10) }
	})
	b.Run("Expensive-n-30", func(b *testing.B) {
		for i := 0; i < b.N; i++ { r = ExpensiveCalculation(30) }
	})

	b.Run("Optimized-n-10", func(b *testing.B) {
		for i := 0; i < b.N; i++ { r = OptimizedCalculation(10) }
	})
	b.Run("Optimized-n-30", func(b *testing.B) {
		for i := 0; i < b.N; i++ { r = OptimizedCalculation(30) }
	})
	b.Run("Optimized-n-1000", func(b *testing.B) {
		for i := 0; i < b.N; i++ { r = OptimizedCalculation(1000) }
	})

	benchResult = r
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
// It should produce identical results but perform better with fewer allocations
func OptimizedSearch(text, substr string) map[int]string {
	if len(substr) > len(text) {
		return map[int]string{}
	}

	if len(substr) == 0 {
		if len(text) == 0 {
			return map[int]string{}
		}
		return map[int]string{0: ""}
	}

	subRunes := make([]rune, 0, len(substr))
	for _, r := range substr {
		subRunes = append(subRunes, unicode.ToLower(r))
	}

	lsp := lspBuild(subRunes)

	maxMatches := len(text) / len(substr)
	if maxMatches == 0 {
		maxMatches = 1
	}
	matches := make(map[int]string, maxMatches)
	
	bytePositions := make([]int, len(subRunes))

	j := 0
	for byteIdx, r := range text {
		lowerR := unicode.ToLower(r)
		
		for j > 0 && lowerR != subRunes[j] {
			oldJ := j
			j = lsp[j-1]
			if j > 0 {
				copy(bytePositions[0:j], bytePositions[oldJ-j:oldJ])
			}
		}
		
		if lowerR == subRunes[j] {
			bytePositions[j] = byteIdx
			j++
		}
		
		if j == len(subRunes) {
			startByteIdx := bytePositions[0]
			endByteIdx := byteIdx + utf8.RuneLen(r)
			
			matches[startByteIdx] = text[startByteIdx:endByteIdx]
			
			oldJ := j
			j = lsp[j-1]
			if j > 0 {
				copy(bytePositions[0:j], bytePositions[oldJ-j:oldJ])
			}
		}
	}

	return matches
}

// lspBuild is helper function for KMP algorithm
func lspBuild(pat []rune) []int {
    lsp := make([]int, len(pat))
    
    j := 0
    for i := 1; i < len(pat); i++ {
        for j > 0 && pat[i] != pat[j] {
            j = lsp[j-1]
        }
        if pat[i] == pat[j] {
            j++
        }
        lsp[i] = j
    }
    return lsp
}

// searchResul defends from compile optimizations on BenchmarkSearchAlgorithms
var searchResult map[int]string

// BenchmarkSearchAlgorithms starts benchmark for Search part of challenge
func BenchmarkSearchAlgorithms(b *testing.B) {
	baseText := "Golang is an open-source language. Sometimes we write GoLaNg, sometimes GOLANG. "
	largeText := strings.Repeat(baseText, 500)
	substr := "golang"

	b.Run("HighAllocation", func(b *testing.B) {
		var r map[int]string
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r = HighAllocationSearch(largeText, substr)
		}
		searchResult = r
	})

	b.Run("OptimizedKMP", func(b *testing.B) {
		var r map[int]string
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r = OptimizedSearch(largeText, substr)
		}
		searchResult = r
	})
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
