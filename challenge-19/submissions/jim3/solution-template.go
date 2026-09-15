package main

import (
	"fmt"
)

// FindMax returns the maximum value in a slice of integers.
// If the slice is empty, it returns 0.
func FindMax(numbers []int) int {
	if len(numbers) == 0 {
		return 0 // or another default value
	}
	max := numbers[0]
	for _, n := range numbers[1:] {
		if n > max {
			max = n
		}
	}

	return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	if len(numbers) == 0 {
		return []int{}
	}
	result := []int{}
	seen := make(map[int]bool)

	for _, v := range numbers {
		_, ok := seen[v]
		if !ok { // Its a new one!
			result = append(result, v)
			seen[v] = true
		}
	}

	return result
}

// ReverseSlice returns a new slice with elements in reverse order
func ReverseSlice(slice []int) []int {
	l := len(slice)
	revSlc := make([]int, l)

	for i, v := range slice {
		revSlc[len(slice)-1-i] = v
	}
	return revSlc
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	result := make([]int, 0)

	for _, n := range numbers {
		if n%2 == 0 {
			result = append(result, n)
		}
	}

	return result
}

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
