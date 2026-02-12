package main

import (
	"slices"
	"strings"
	"time"
)

func SlowSort(data []int) []int {
	result := make([]int, len(data))
	copy(result, data)
	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result)-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}

func OptimizedSort(data []int) []int {
	result := make([]int, len(data))
	copy(result, data)
	slices.Sort(result)
	return result
}

func InefficientStringBuilder(parts []string, repeatCount int) string {
	var buff = strings.Builder{}
	for i := 0; i < repeatCount; i++ {
		for _, part := range parts {
			buff.WriteString(part)
		}
	}
	return buff.String()
}

func OptimizedStringBuilder(parts []string, repeatCount int) string {
	var combined strings.Builder
	for _, part := range parts {
		combined.WriteString(part)
	}
	combinedStr := combined.String()

	var result strings.Builder
	result.Grow(len(combinedStr) * repeatCount)
	for i := 0; i < repeatCount; i++ {
		result.WriteString(combinedStr)
	}
	return result.String()
}

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

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func OptimizedCalculation(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	a, b := 0, 1
	sum := 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
		sum += b
	}
	return sum
}

func HighAllocationSearch(text, substr string) map[int]string {
	result := make(map[int]string)
	lowerText := strings.ToLower(text)
	lowerSubstr := strings.ToLower(substr)
	for i := 0; i < len(lowerText); i++ {
		if i+len(lowerSubstr) <= len(lowerText) {
			potentialMatch := lowerText[i : i+len(lowerSubstr)]
			if potentialMatch == lowerSubstr {
				result[i] = text[i : i+len(substr)]
			}
		}
	}
	return result
}

func OptimizedSearch(text, substr string) map[int]string {
	result := make(map[int]string)
	if substr == "" {
		return result
	}
	
	lowerSubstr := strings.ToLower(substr)
	lowerText := strings.ToLower(text)
	
	start := 0
	for start < len(lowerText) {
		idx := strings.Index(lowerText[start:], lowerSubstr)
		if idx == -1 {
			break
		}
		pos := start + idx
		result[pos] = text[pos : pos+len(substr)]
		start = pos + 1
	}
	return result
}

func SimulateCPUWork(duration time.Duration) {
	start := time.Now()
	for time.Since(start) < duration {
		for i := 0; i < 1000000; i++ {
			_ = i
		}
	}
}