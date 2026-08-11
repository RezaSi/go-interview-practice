package main

import (
	"fmt"
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
	var result = numbers[0]
	for _, number := range numbers {
		if number > result {
			result = number
		}
	}
	return result
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	result := make([]int, 0, len(numbers))
	corr := make(map[int]struct{}, len(numbers))
	for _, number := range numbers {
		if _, found := corr[number]; !found {
			result = append(result, number)
			corr[number] = struct{}{}
		}
	}
	return result
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	result := make([]int, 0, len(slice))
	for i := len(slice) - 1; i >= 0; i-- {
		result = append(result, slice[i])
	}
	return result
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	result := make([]int, 0, len(numbers))
	for _, number := range numbers {
		if number%2 == 0 {
			result = append(result, number)
		}
	}
	return result
}
