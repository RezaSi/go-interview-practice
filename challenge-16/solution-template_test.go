package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

// Helper function to generate random slices for testing
func generateRandomSlice(size int) []int {
	rand.Seed(time.Now().UnixNano())
	slice := make([]int, size)
	for i := 0; i < size; i++ {
		slice[i] = rand.Intn(10000)
	}
	return slice
}

// Helper function to check if two slices are equal
func slicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestSlowSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
	}{
		{"Empty", []int{}},
		{"One Element", []int{42}},
		{"Already Sorted", []int{1, 2, 3, 4, 5}},
		{"Reverse Sorted", []int{5, 4, 3, 2, 1}},
		{"Random Order", []int{3, 1, 4, 1, 5, 9, 2, 6}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a copy and sort it using Go's built-in sort for comparison
			expected := make([]int, len(tc.input))
			copy(expected, tc.input)
			sort.Ints(expected)

			// Test our slow sort
			result := SlowSort(tc.input)
			if !slicesEqual(result, expected) {
				t.Errorf("SlowSort didn't sort correctly. Got %v, expected %v", result, expected)
			}
		})
	}
}

func TestOptimizedSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
	}{
		{"Empty", []int{}},
		{"One Element", []int{42}},
		{"Already Sorted", []int{1, 2, 3, 4, 5}},
		{"Reverse Sorted", []int{5, 4, 3, 2, 1}},
		{"Random Order", []int{3, 1, 4, 1, 5, 9, 2, 6}},
		{"Random Large", generateRandomSlice(100)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get results from both methods
			slowResult := SlowSort(tc.input)
			optimizedResult := OptimizedSort(tc.input)

			// Check that optimized sort gives the same result as slow sort
			if !slicesEqual(slowResult, optimizedResult) {
				t.Errorf("OptimizedSort gave different results than SlowSort. Got %v, expected %v", optimizedResult, slowResult)
			}
		})
	}
}

func BenchmarkSlowSort(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		data := generateRandomSlice(size)
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				SlowSort(data)
			}
		})
	}
}

func BenchmarkOptimizedSort(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		data := generateRandomSlice(size)
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				OptimizedSort(data)
			}
		})
	}
}

func TestInefficientStringBuilder(t *testing.T) {
	testCases := []struct {
		name        string
		parts       []string
		repeatCount int
	}{
		{"Empty", []string{}, 10},
		{"Single Part", []string{"Hello"}, 5},
		{"Multiple Parts", []string{"Hello", " ", "World", "!"}, 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := InefficientStringBuilder(tc.parts, tc.repeatCount)

			// Calculate expected result
			expected := ""
			for i := 0; i < tc.repeatCount; i++ {
				for _, part := range tc.parts {
					expected += part
				}
			}

			if result != expected {
				t.Errorf("InefficientStringBuilder result incorrect. Got %q, expected %q", result, expected)
			}
		})
	}
}

func TestOptimizedStringBuilder(t *testing.T) {
	testCases := []struct {
		name        string
		parts       []string
		repeatCount int
	}{
		{"Empty", []string{}, 10},
		{"Single Part", []string{"Hello"}, 5},
		{"Multiple Parts", []string{"Hello", " ", "World", "!"}, 3},
		{"Large Build", []string{"This", " ", "is", " ", "a", " ", "test"}, 100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inefficientResult := InefficientStringBuilder(tc.parts, tc.repeatCount)
			optimizedResult := OptimizedStringBuilder(tc.parts, tc.repeatCount)

			if optimizedResult != inefficientResult {
				t.Errorf("OptimizedStringBuilder gave different results than InefficientStringBuilder. Got %q, expected %q", optimizedResult, inefficientResult)
			}
		})
	}
}

func BenchmarkInefficientStringBuilder(b *testing.B) {
	testCases := []struct {
		name        string
		parts       []string
		repeatCount int
	}{
		{"Small", []string{"Hello", " ", "World"}, 10},
		{"Medium", []string{"This", " ", "is", " ", "a", " ", "test"}, 100},
		{"Large", []string{"The", " ", "quick", " ", "brown", " ", "fox"}, 1000},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				InefficientStringBuilder(tc.parts, tc.repeatCount)
			}
		})
	}
}

func BenchmarkOptimizedStringBuilder(b *testing.B) {
	testCases := []struct {
		name        string
		parts       []string
		repeatCount int
	}{
		{"Small", []string{"Hello", " ", "World"}, 10},
		{"Medium", []string{"This", " ", "is", " ", "a", " ", "test"}, 100},
		{"Large", []string{"The", " ", "quick", " ", "brown", " ", "fox"}, 1000},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				OptimizedStringBuilder(tc.parts, tc.repeatCount)
			}
		})
	}
}

func TestExpensiveCalculation(t *testing.T) {
	testCases := []struct {
		n        int
		expected int
	}{
		{0, 0},
		{1, 1},
		{5, 12},
		{10, 143},
	}

	for _, tc := range testCases {
		result := ExpensiveCalculation(tc.n)
		if result != tc.expected {
			t.Errorf("ExpensiveCalculation(%d) = %d, expected %d", tc.n, result, tc.expected)
		}
	}
}

func TestOptimizedCalculation(t *testing.T) {
	testCases := []struct {
		n int
	}{
		{0},
		{1},
		{5},
		{10},
		{15},
	}

	for _, tc := range testCases {
		expensiveResult := ExpensiveCalculation(tc.n)
		optimizedResult := OptimizedCalculation(tc.n)

		if optimizedResult != expensiveResult {
			t.Errorf("OptimizedCalculation(%d) = %d, expected %d", tc.n, optimizedResult, expensiveResult)
		}
	}
}

func BenchmarkExpensiveCalculation(b *testing.B) {
	benchmarks := []struct {
		name string
		n    int
	}{
		{"Small", 10},
		{"Medium", 20},
		{"Large", 30},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ExpensiveCalculation(bm.n)
			}
		})
	}
}

func BenchmarkOptimizedCalculation(b *testing.B) {
	benchmarks := []struct {
		name string
		n    int
	}{
		{"Small", 10},
		{"Medium", 20},
		{"Large", 30},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				OptimizedCalculation(bm.n)
			}
		})
	}
}

func TestHighAllocationSearch(t *testing.T) {
	testCases := []struct {
		name    string
		text    string
		substr  string
		matches int
	}{
		{"No Match", "Hello World", "xyz", 0},
		{"Single Match", "Hello World", "world", 1},
		{"Multiple Matches", "banana", "an", 2},
		{"Case Insensitive", "Hello World Hello", "hello", 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := HighAllocationSearch(tc.text, tc.substr)
			if len(result) != tc.matches {
				t.Errorf("Expected %d matches, got %d", tc.matches, len(result))
			}
		})
	}
}

func TestOptimizedSearch(t *testing.T) {
	testCases := []struct {
		name   string
		text   string
		substr string
	}{
		{"Empty", "", ""},
		{"No Match", "Hello World", "xyz"},
		{"Single Match", "Hello World", "world"},
		{"Multiple Matches", "banana", "an"},
		{"Case Insensitive", "Hello World Hello", "hello"},
		{"Long Text", strings.Repeat("The quick brown fox jumps over the lazy dog. ", 100), "fox"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inefficientResult := HighAllocationSearch(tc.text, tc.substr)
			optimizedResult := OptimizedSearch(tc.text, tc.substr)

			if !reflect.DeepEqual(inefficientResult, optimizedResult) {
				t.Errorf("OptimizedSearch gave different results than HighAllocationSearch")
			}
		})
	}
}

func BenchmarkHighAllocationSearch(b *testing.B) {
	benchmarks := []struct {
		name   string
		text   string
		substr string
	}{
		{"Short Text", "The quick brown fox jumps over the lazy dog.", "fox"},
		{"Medium Text", strings.Repeat("The quick brown fox jumps over the lazy dog. ", 10), "fox"},
		{"Long Text", strings.Repeat("The quick brown fox jumps over the lazy dog. ", 100), "fox"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				HighAllocationSearch(bm.text, bm.substr)
			}
		})
	}
}

func BenchmarkOptimizedSearch(b *testing.B) {
	benchmarks := []struct {
		name   string
		text   string
		substr string
	}{
		{"Short Text", "The quick brown fox jumps over the lazy dog.", "fox"},
		{"Medium Text", strings.Repeat("The quick brown fox jumps over the lazy dog. ", 10), "fox"},
		{"Long Text", strings.Repeat("The quick brown fox jumps over the lazy dog. ", 100), "fox"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				OptimizedSearch(bm.text, bm.substr)
			}
		})
	}
}

func BenchmarkMemoryHighAllocationSearch(b *testing.B) {
	text := strings.Repeat("The quick brown fox jumps over the lazy dog. ", 100)
	substr := "fox"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		HighAllocationSearch(text, substr)
	}
}

func BenchmarkMemoryOptimizedSearch(b *testing.B) {
	text := strings.Repeat("The quick brown fox jumps over the lazy dog. ", 100)
	substr := "fox"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		OptimizedSearch(text, substr)
	}
}
