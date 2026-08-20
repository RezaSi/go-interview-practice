package main

import (
	"fmt"
	"math"
)

func main() {
	// Example slice for testing
	numbers := []int{3, 1, 4, 1, 5, 9, 2, 6}

	// Test FindMax
	max := FindMax(numbers)
	fmt.Printf("Maximum value: %d\n", max)

	// Test RemoveDuplicates
	unique := RemoveDuplicates(numbers)
	fmt.Printf("After removing duplicates: %v\n", unique)

	// Test ReverseSlice
	reversed := ReverseSlice(numbers)
	fmt.Printf("Reversed: %v\n", reversed)

	// Test FilterEven
	evenOnly := FilterEven(numbers)
	fmt.Printf("Even numbers only: %v\n", evenOnly)
}
// FindMax returns the maximum value in a slice of integers.
// If the slice is empty, it returns 0.
func FindMax(numbers []int) int {
	if len(numbers) == 0 {
		return 0
	}
	maxElement := math.MinInt
	for _, x := range numbers {
		maxElement = max(maxElement, x)
	}

	return maxElement
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	if len(numbers) == 0 {
		return numbers
	}
	marked := make(map[int]struct{})
	var result []int
	for _, x := range numbers {
		if _, ok := marked[x]; !ok {
			marked[x] = struct{}{}
			result = append(result, x)
		}
	}
	return result
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	i := 0
	j := len(slice) - 1
	result := make([]int, len(slice))
	for i <= j {
		result[i], result[j] = slice[j], slice[i]
		i++
		j--
	}
	return result
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	result := []int{}
	for _, x := range numbers {
		if x%2 == 0 {
			result = append(result, x)
		}
	}
	return result
}

