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

	maxNum := numbers[0]
	for _, num := range numbers {
		maxNum = max(maxNum, num)
	}

	return maxNum
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	seen := make(map[int]bool)
	unique := []int{}

	for _, num := range numbers {
		if !seen[num] {
			seen[num] = true
			unique = append(unique, num)
		}
	}

	return unique
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	reversedSlice := make([]int, 0)
	for i := len(slice) - 1; i >= 0; i-- {
		reversedSlice = append(reversedSlice, slice[i])
	}

	return reversedSlice
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	evenOnly := []int{}

	for _, num := range numbers {
		if num%2 == 0 {
			evenOnly = append(evenOnly, num)
		}
	}

	return evenOnly
}
