package main

import (
	"slices"
	"strings"
	"time"
	"unicode"
)

// Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkSlowSort$ challenge16 -count=1

// goos: darwin
// goarch: arm64
// pkg: challenge16
// cpu: Apple M1 Pro
// === RUN   BenchmarkSlowSort
// BenchmarkSlowSort
// === RUN   BenchmarkSlowSort/10
// BenchmarkSlowSort/10
// BenchmarkSlowSort/10-10         16799235                70.37 ns/op           80 B/op          1 allocs/op
// === RUN   BenchmarkSlowSort/100
// BenchmarkSlowSort/100
// BenchmarkSlowSort/100-10                  176282              6763 ns/op             896 B/op          1 allocs/op
// === RUN   BenchmarkSlowSort/1000
// BenchmarkSlowSort/1000
// BenchmarkSlowSort/1000-10                   1690            714555 ns/op            8192 B/op          1 allocs/op
// PASS
// ok      challenge16     5.019s

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

// Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkOptimizedSort$ challenge16 -count=1

// goos: darwin
// goarch: arm64
// pkg: challenge16
// cpu: Apple M1 Pro
// === RUN   BenchmarkOptimizedSort
// BenchmarkOptimizedSort
// === RUN   BenchmarkOptimizedSort/10
// BenchmarkOptimizedSort/10
// BenchmarkOptimizedSort/10-10    115245958               10.24 ns/op            0 B/op          0 allocs/op
// === RUN   BenchmarkOptimizedSort/100
// BenchmarkOptimizedSort/100
// BenchmarkOptimizedSort/100-10           15166443                78.99 ns/op            0 B/op          0 allocs/op
// === RUN   BenchmarkOptimizedSort/1000
// BenchmarkOptimizedSort/1000
// BenchmarkOptimizedSort/1000-10           1826474               683.1 ns/op             0 B/op          0 allocs/op
// PASS

// OptimizedSort is your optimized version of SlowSort
// It should produce identical results but perform better
func OptimizedSort(data []int) []int {
	slices.Sort(data)
	return data
}

// Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkInefficientStringBuilder$ challenge16 -count=1

// goos: darwin
// goarch: arm64
// pkg: challenge16
// cpu: Apple M1 Pro
// === RUN   BenchmarkInefficientStringBuilder
// BenchmarkInefficientStringBuilder
// === RUN   BenchmarkInefficientStringBuilder/Small
// BenchmarkInefficientStringBuilder/Small
// BenchmarkInefficientStringBuilder/Small-10               1501484               765.3 ns/op          1912 B/op         29 allocs/op
// === RUN   BenchmarkInefficientStringBuilder/Medium
// BenchmarkInefficientStringBuilder/Medium
// BenchmarkInefficientStringBuilder/Medium-10                18094             62895 ns/op          518165 B/op        699 allocs/op
// === RUN   BenchmarkInefficientStringBuilder/Large
// BenchmarkInefficientStringBuilder/Large
// BenchmarkInefficientStringBuilder/Large-10                   190           5438039 ns/op        70153518 B/op       7001 allocs/op
// PASS

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

// Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkOptimizedStringBuilder$ challenge16 -count=1

// goos: darwin
// goarch: arm64
// pkg: challenge16
// cpu: Apple M1 Pro
// === RUN   BenchmarkOptimizedStringBuilder
// BenchmarkOptimizedStringBuilder
// === RUN   BenchmarkOptimizedStringBuilder/Small
// BenchmarkOptimizedStringBuilder/Small
// BenchmarkOptimizedStringBuilder/Small-10                16017208                70.95 ns/op          128 B/op          2 allocs/op
// === RUN   BenchmarkOptimizedStringBuilder/Medium
// BenchmarkOptimizedStringBuilder/Medium
// BenchmarkOptimizedStringBuilder/Medium-10                5503075               204.6 ns/op          1424 B/op          2 allocs/op
// === RUN   BenchmarkOptimizedStringBuilder/Large
// BenchmarkOptimizedStringBuilder/Large
// BenchmarkOptimizedStringBuilder/Large-10                  876415              1339 ns/op           19096 B/op          2 allocs/op
// PASS

// OptimizedStringBuilder is your optimized version of InefficientStringBuilder
// It should produce identical results but perform better
func OptimizedStringBuilder(parts []string, repeatCount int) string {
	// Could be better if using a loop with strings.Builder, but this solution is shorter, and performs ok.
	return strings.Repeat(strings.Join(parts, ""), repeatCount)
}

// Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkExpensiveCalculation$ challenge16 -count=1

// goos: darwin
// goarch: arm64
// pkg: challenge16
// cpu: Apple M1 Pro
// === RUN   BenchmarkExpensiveCalculation
// BenchmarkExpensiveCalculation
// === RUN   BenchmarkExpensiveCalculation/Small
// BenchmarkExpensiveCalculation/Small
// BenchmarkExpensiveCalculation/Small-10           2502798               487.3 ns/op             0 B/op          0 allocs/op
// === RUN   BenchmarkExpensiveCalculation/Medium
// BenchmarkExpensiveCalculation/Medium
// BenchmarkExpensiveCalculation/Medium-10            20091             59268 ns/op               0 B/op          0 allocs/op
// === RUN   BenchmarkExpensiveCalculation/Large
// BenchmarkExpensiveCalculation/Large
// BenchmarkExpensiveCalculation/Large-10               163           7262633 ns/op               0 B/op          0 allocs/op
// PASS

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

// Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkOptimizedCalculation$ challenge16 -count=1

// goos: darwin
// goarch: arm64
// pkg: challenge16
// cpu: Apple M1 Pro
// === RUN   BenchmarkOptimizedCalculation
// BenchmarkOptimizedCalculation
// === RUN   BenchmarkOptimizedCalculation/Small
// BenchmarkOptimizedCalculation/Small
// BenchmarkOptimizedCalculation/Small-10          92842292                12.66 ns/op            0 B/op          0 allocs/op
// === RUN   BenchmarkOptimizedCalculation/Medium
// BenchmarkOptimizedCalculation/Medium
// BenchmarkOptimizedCalculation/Medium-10         54438536                22.03 ns/op            0 B/op          0 allocs/op
// === RUN   BenchmarkOptimizedCalculation/Large
// BenchmarkOptimizedCalculation/Large
// BenchmarkOptimizedCalculation/Large-10          37801525                31.78 ns/op            0 B/op          0 allocs/op
// PASS

// OptimizedCalculation is your optimized version of ExpensiveCalculation
// It should produce identical results but perform better
func OptimizedCalculation(n int) int {
	// sum = f(n+2)-1
	return fastFib(n+2) - 1
}

func fastFib(n int) int {
	return fib2(1, 1, n)
}

func fib2(a, b, n int) int {
	if n <= 1 {
		return a
	}
	return fib2(b, a+b, n-1)
}

// Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkHighAllocationSearch$ challenge16 -count=1

// goos: darwin
// goarch: arm64
// pkg: challenge16
// cpu: Apple M1 Pro
// === RUN   BenchmarkHighAllocationSearch
// BenchmarkHighAllocationSearch
// === RUN   BenchmarkHighAllocationSearch/Short_Text
// BenchmarkHighAllocationSearch/Short_Text
// BenchmarkHighAllocationSearch/Short_Text-10              4302798               294.2 ns/op           304 B/op          3 allocs/op
// === RUN   BenchmarkHighAllocationSearch/Medium_Text
// BenchmarkHighAllocationSearch/Medium_Text
// BenchmarkHighAllocationSearch/Medium_Text-10              513793              2324 ns/op            1192 B/op          6 allocs/op
// === RUN   BenchmarkHighAllocationSearch/Long_Text
// BenchmarkHighAllocationSearch/Long_Text
// BenchmarkHighAllocationSearch/Long_Text-10                 53697             22557 ns/op           11816 B/op         12 allocs/op
// PASS

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

// Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkOptimizedSearch$ challenge16 -count=1

// goos: darwin
// goarch: arm64
// pkg: challenge16
// cpu: Apple M1 Pro
// === RUN   BenchmarkOptimizedSearch
// BenchmarkOptimizedSearch
// === RUN   BenchmarkOptimizedSearch/Short_Text
// BenchmarkOptimizedSearch/Short_Text
// BenchmarkOptimizedSearch/Short_Text-10           4806428               235.2 ns/op           256 B/op          2 allocs/op
// === RUN   BenchmarkOptimizedSearch/Medium_Text
// BenchmarkOptimizedSearch/Medium_Text
// BenchmarkOptimizedSearch/Medium_Text-10           592110              1997 ns/op             712 B/op          5 allocs/op
// === RUN   BenchmarkOptimizedSearch/Long_Text
// BenchmarkOptimizedSearch/Long_Text
// BenchmarkOptimizedSearch/Long_Text-10              59658             20095 ns/op            6952 B/op         11 allocs/op
// PASS

// OptimizedSearch is your optimized version of HighAllocationSearch
// It should produce identical results but perform better with fewer allocations
func OptimizedSearch(text, substr string) map[int]string {
	if len(text) == 0 {
		return map[int]string{}
	}
	res := map[int]string{}
	for i := 0; i <= len(text)-len(substr); i++ {
		if isMatch(text, substr, i) {
			res[i] = text[i : i+len(substr)]
		}
	}
	return res
}

func isMatch(text, substr string, begin int) bool {
	for i := 0; i < len(substr); i++ {
		if unicode.ToLower(rune(text[begin+i])) != unicode.ToLower(rune(substr[i])) {
			return false
		}
	}
	return true
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
