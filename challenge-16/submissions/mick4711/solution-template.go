package main

import (
	// "bytes"
	// index/suffixarray
	"slices"
	"strings"
	"time"
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
	// Hint: Consider using sort package or a more efficient algorithm

	// sort.Slice(data, func(i, j int) bool {
	// 	return data[i] < data[j]
	// })

	slices.Sort(data)

	/*
			   go test -v -timeout 30s -run="SlowSort|OptimizedSort"
			   go test -benchmem -bench="SlowSort|OptimizedSort"

			   BenchmarkSlowSort/10-8           8783551             123.9 ns/op            80 B/op          1 allocs/op
			   BenchmarkSlowSort/100-8           107521             11458 ns/op           896 B/op          1 allocs/op
			   BenchmarkSlowSort/1000-8             607           1653977 ns/op          8192 B/op          1 allocs/op
		sort.Slice
			   BenchmarkOptimizedSort/10-8      9244336             159.2 ns/op            56 B/op          2 allocs/op
			   BenchmarkOptimizedSort/100-8     2408090             532.3 ns/op            56 B/op          2 allocs/op
			   BenchmarkOptimizedSort/1000-8     313731              3719 ns/op            56 B/op          2 allocs/op
		slices.Sort
			   BenchmarkOptimizedSort/10-8     63324622             18.60 ns/op             0 B/op          0 allocs/op
			   BenchmarkOptimizedSort/100-8     9154833             133.4 ns/op             0 B/op          0 allocs/op
			   BenchmarkOptimizedSort/1000-8    1000000              1162 ns/op             0 B/op          0 allocs/op
	*/

	return data
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
	// Hint: Consider using strings.Builder or bytes.Buffer

	// var b bytes.Buffer // A Buffer needs no initialization.
	var b strings.Builder

	for range repeatCount {
		for _, part := range parts {
			b.Write([]byte(part))
		}
	}

	/*
			go test -v -timeout 30s -run="StringBuilder"
			go test -benchmem -bench="StringBuilder"

			BenchmarkInefficientStringBuilder/Small-8      836750         1385 ns/op        1912 B/op         29 allocs/op
			BenchmarkInefficientStringBuilder/Medium-8       6750       187868 ns/op      518168 B/op        699 allocs/op
			BenchmarkInefficientStringBuilder/Large-8          73     17936059 ns/op    70153736 B/op       7003 allocs/op
		bytes.Buffer
			BenchmarkOptimizedStringBuilder/Small-8       3746812          306.2 ns/op       304 B/op          3 allocs/op
			BenchmarkOptimizedStringBuilder/Medium-8       228163         4949 ns/op        5440 B/op          7 allocs/op
			BenchmarkOptimizedStringBuilder/Large-8         22855        54112 ns/op       84544 B/op         11 allocs/op
		strings.Builder
			BenchmarkOptimizedStringBuilder/Small-8       4266385          289.4 ns/op       248 B/op          5 allocs/op
			BenchmarkOptimizedStringBuilder/Medium-8       344193         3603 ns/op        3320 B/op          9 allocs/op
			BenchmarkOptimizedStringBuilder/Large-8         25576        47915 ns/op       84728 B/op         18 allocs/op
	*/

	return b.String()
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
	// Hint: Consider memoization or avoiding redundant calculations
	/*
			go test -v -timeout 30s -run="Calculation"
			go test -benchmem -bench="Calculation"

			BenchmarkExpensiveCalculation/Small-8      2035190         595.9 ns/op     0 B/op     0 allocs/op
			BenchmarkExpensiveCalculation/Medium-8       16122       75309 ns/op       0 B/op     0 allocs/op
			BenchmarkExpensiveCalculation/Large-8          127     9544298 ns/op       0 B/op     0 allocs/op
		Memoize
			BenchmarkOptimizedCalculation/Small-8      2366942       498.8 ns/op     328 B/op     3 allocs/op
			BenchmarkOptimizedCalculation/Medium-8     1256305       949.1 ns/op     616 B/op     3 allocs/op
			BenchmarkOptimizedCalculation/Large-8       755708      1451 ns/op      1192 B/op     3 allocs/op
		Iterate
			BenchmarkOptimizedCalculation/Small-8    340811240         3.555 ns/op     0 B/op     0 allocs/op
			BenchmarkOptimizedCalculation/Medium-8   201614457         6.047 ns/op     0 B/op     0 allocs/op
			BenchmarkOptimizedCalculation/Large-8    144884427         8.462 ns/op     0 B/op     0 allocs/op
	*/
	// return ExpensiveCalculation(n) // Replace this with your optimized implementation
	if n <= 1 {
		return n
	}

	sum := 0

	// Memoize
	/*
		memo := make(map[int]int, n)

		for i := 1; i <= n; i++ {
			sum += fibonacciMem(i, memo)
		}
	*/

	// Iterate
	prev1, prev2 := 1, 0
	for i := 1; i <= n; i++ {
		prev2, prev1 = prev1, prev2+prev1
		sum += prev2
	}

	return sum
}

// Helper function that computes the fibonacci number at position n using memoization
func fibonacciMem(n int, memo map[int]int) int {
	if n <= 1 {
		return n
	}

	res1, ok := memo[n-1]
	if !ok {
		res1 = fibonacciMem(n-1, memo)
		memo[n-1] = res1
	}

	res2, ok := memo[n-2]
	if !ok {
		res2 = fibonacciMem(n-2, memo)
		memo[n-2] = res2
	}

	return res1 + res2
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
	// Hint: Consider avoiding temporary string allocations and reusing memory
	/*
			go test -v -timeout 30s -run="Search"
			go test -benchmem -bench="Search"

			BenchmarkHighAllocationSearch/Short_Text-8     3013060        393.6 ns/op      304 B/op     3 allocs/op
			BenchmarkHighAllocationSearch/Medium_Text-8     370346       3083 ns/op       1192 B/op     6 allocs/op
			BenchmarkHighAllocationSearch/Long_Text-8        40744      29822 ns/op      11816 B/op    12 allocs/op
		index/suffixarray
			BenchmarkOptimizedSearch/Short_Text-8           577708       1891 ns/op        616 B/op     7 allocs/op
			BenchmarkOptimizedSearch/Medium_Text-8           91478      13193 ns/op       3880 B/op    10 allocs/op
			BenchmarkOptimizedSearch/Long_Text-8              9838     102849 ns/op      36088 B/op    16 allocs/op
		strings.Index
			BenchmarkOptimizedSearch/Short_Text-8          5038852        237.9 ns/op      304 B/op     3 allocs/op
			BenchmarkOptimizedSearch/Medium_Text-8          786710       1397 ns/op       1192 B/op     6 allocs/op
			BenchmarkOptimizedSearch/Long_Text-8             94105      12672 ns/op      11816 B/op    12 allocs/op
	*/
	result := make(map[int]string)

	// Convert to lowercase for case-insensitive search
	lowerText := strings.ToLower(text)
	lowerSubstr := strings.ToLower(substr)

	lt := len(lowerText)
	ls := len(lowerSubstr)
	left := 0
	for left < lt {
		i := strings.Index(lowerText[left:], lowerSubstr)
		if i < 0 {
			break
		}
		left += i
		result[left] = text[left : left+ls]
		left += ls
	}

	/*
		ls := len(lowerSubstr)
		// create index for some data
		index := suffixarray.New([]byte(lowerText))

		// lookup byte slice s
		offsets1 := index.Lookup([]byte(lowerSubstr), -1) // the list of all indices where s occurs in data
		for _, i := range offsets1 {
			result[i] = text[i : i+ls]
		}
	*/
	/*
		lt := len(lowerText)
		ls := len(lowerSubstr)
		for i := 0; i < lt; i++ {
			// Check if we can fit the substring starting at position i
			if i+ls <= lt {
				// Extract the potential match
				potentialMatch := lowerText[i : i+ls]

				// Check if it matches
				if potentialMatch == lowerSubstr {
					// Store the original case version
					result[i] = text[i : i+ls]
				}
			}
		}
	*/

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
