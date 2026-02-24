package main

import (
	"sort"
	"strings"
	"time"
)

func SlowSort(data []int) []int {
	result := make([]int, len(data))
	copy(result, data)
	for i := 0; i < len(result)-1; i++ {
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
	sort.Ints(result)
	return result
}

func InefficientStringBuilder(parts []string, repeatCount int) string {
	result := ""
	for i := 0; i < repeatCount; i++ {
		for _, part := range parts {
			result += part
		}
	}
	return result
}
func OptimizedStringBuilder(words []string, repeatCount int) string {
	var builder strings.Builder
	totalLen := 0
	for _, word := range words {
		totalLen += len(word)
	}
	builder.Grow(totalLen * repeatCount)
	for range repeatCount {
		for _, word := range words {
			builder.WriteString(word)
		}
	}
	return builder.String()
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
	if n <= 2 {
		return n
	}
	f1, f2 := 1, 1
	sum := 2
	for range n - 2 {
		f1, f2 = f2, f1+f2
		sum += f2
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
	if text == "" || substr == "" {
		return map[int]string{}
	}
	result := make(map[int]string)
	for i := 0; i <= len(text)-len(substr); i++ {
		candidate := text[i : i+len(substr)]
		if strings.EqualFold(candidate, substr) {
			result[i] = text[i : i+len(substr)]
		}
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