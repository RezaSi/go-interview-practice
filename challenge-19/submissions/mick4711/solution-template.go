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
	// check for empty slice
	if len(numbers) == 0 {
		return 0
	}

	// init with first value
	max := numbers[0]

	// check subsequent vals for greater value
	for _, v := range numbers {
		if v > max {
			max = v
		}
	}

	return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	// allocate result slice
	res := make([]int, 0, len(numbers))

	// loop thru slice filtering out dups
	for _, v := range numbers {
		if !slices.Contains(res, v) {
			res = append(res, v)
		}
	}

	return res
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	// allocate result slice
	l := len(slice)
	res := make([]int, l)

	// loop thru slice filling res from the end backwards
	for i, v := range slice {
		res[l-i-1] = v
	}

	return res
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	// allocate result slice
	res := make([]int, 0, len(numbers))

	// loop thru slice filtering out dups
	for _, v := range numbers {
		if v%2 == 0 {
			res = append(res, v)
		}
	}

	return res
}
