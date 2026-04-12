package main

import (
	"fmt"
	"slices"
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
	return slices.Max(numbers)
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	seen := map[int]bool{}
	deduped := []int{}
	for _, n := range numbers {
		if seen[n] {
			continue
		}
		seen[n] = true
		deduped = append(deduped, n)
	}
	return deduped
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	res := make([]int, 0, len(slice))
	for i := len(slice) - 1; i >= 0; i-- {
		res = append(res, slice[i])
	}
	return res
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	res := []int{}
	for _, n := range numbers {
		if n%2 == 0 {
			res = append(res, n)
		}
	}
	return res
}
