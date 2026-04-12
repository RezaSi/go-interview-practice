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
	
	maximum := numbers[0]
	for i := 1; i < len(numbers); i++ {
		maximum = max(maximum, numbers[i])
	}
	
	return maximum
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	result := make([]int, 0)
	m := make(map[int]bool)
	
	for _, val := range numbers {
		if !m[val] {
			m[val] = true
			result = append(result, val)
		}
	}
	
	return result
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	result := make([]int, len(slice))
	
	for i := 0; i < len(slice); i++ {
		result[i] = slice[len(slice)-i-1]
	}
	
	return result
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	result := make([]int, 0)
	
	for _, val := range numbers {
		if val % 2 == 0 {
			result = append(result, val)
		}
	}
	
	return result
}
