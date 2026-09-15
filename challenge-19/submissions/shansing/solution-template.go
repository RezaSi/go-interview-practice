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
	m := numbers[0]
	for _, number := range numbers {
		if number > m {
			m = number
		}
	}
	return m
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	result := make([]int, len(numbers))
	if len(numbers) == 0 {
		return result
	}
	showed := make(map[int]struct{}, len(numbers))
	i := 0
	for _, number := range numbers {
		if _, ok := showed[number]; ok {
			continue
		}
		showed[number] = struct{}{}
		result[i] = number
		i++
	}
	return result[:i]
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	result := make([]int, len(slice))
	for i := 0; i < len(slice); i++ {
		j := len(slice) - 1 - i
		result[j] = slice[i]
	}
	return result
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	result := make([]int, len(numbers))
	i := 0
	for _, number := range numbers {
		if number%2 == 0 {
			result[i] = number
			i++
		}
	}
	return result[:i]
}
