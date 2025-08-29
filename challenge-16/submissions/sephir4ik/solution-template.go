package main

import (
	"strings"
	"time"
	"sort"
	"unicode"
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
	builder := strings.Builder{}
	
	for _, part := range parts {
	    builder.WriteString(part)
	}
	
	str := builder.String()
	
	for i:=1; i < repeatCount; i++{
	    builder.WriteString(str)
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
    if n <= 0 {
        return 0
    }

    sum := 0
    a, b := 0, 1
    
    for i := 1; i <= n; i++ {
        sum += b
        a, b = b, a+b
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

    // Преобразуем подстроку в нижний регистр один раз
    lowerSubstr := strings.ToLower(substr)
    substrRunes := []rune(lowerSubstr)
    substrLen := len(substrRunes)

    // Преобразуем текст в руны один раз с сохранением позиций в байтах
    var textRunes []rune
    var byteOffsets []int
    for pos, r := range text {
        textRunes = append(textRunes, unicode.ToLower(r))
        byteOffsets = append(byteOffsets, pos)
    }
    byteOffsets = append(byteOffsets, len(text))

    // Ищем совпадения
    for i := 0; i <= len(textRunes)-substrLen; i++ {
        match := true
        for j := 0; j < substrLen; j++ {
            if textRunes[i+j] != substrRunes[j] {
                match = false
                break
            }
        }
        if match {
            start := byteOffsets[i]
            end := byteOffsets[i+substrLen]
            result[start] = text[start:end]
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
