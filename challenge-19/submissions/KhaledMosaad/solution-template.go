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
	ret := slices.Max(numbers)
	return ret
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	numLen := len(numbers)
	if numLen == 0 {
		return numbers
	}

	ret := make([]int, 0, numLen)
	uniqueness := make(map[int]bool)

	for i := 0; i < numLen; i++ {
		_, ok := uniqueness[numbers[i]]
		if !ok {
			uniqueness[numbers[i]] = true
			ret = append(ret, numbers[i])
		}
	}
	return ret
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	ret := slices.Clone(slice)
	slices.Reverse(ret)
	return ret
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	numLen := len(numbers)

	if numLen == 0 {
		return numbers
	}

	ret := make([]int, 0, numLen)
	for i := 0; i < numLen; i++ {
		if numbers[i]%2 == 0 {
			ret = append(ret, numbers[i])
		}
	}
	return ret
}
