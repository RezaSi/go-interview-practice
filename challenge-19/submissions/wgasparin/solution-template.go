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

	max := math.MinInt

	for _, v := range numbers {
		if max <= v {
			max = v
		}
	}

	return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	deduplicate := []int{}

	for _, v := range numbers {

		found := false

		for _, v1 := range deduplicate {
			if v == v1 {
				found = true
				break
			}
		}

		if !found || len(deduplicate) == 0 {
			deduplicate = append(deduplicate, v)
		}
	}

	return deduplicate
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	newSlice := []int{}

	for i := len(slice) - 1; i >= 0; i-- {
		newSlice = append(newSlice, slice[i])
	}
	return newSlice
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {

	new := []int{}

	for _, v := range numbers {
		if v%2 == 0 {
			new = append(new, v)
		}
	}

	return new
}
