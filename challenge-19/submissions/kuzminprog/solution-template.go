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
	for _, num := range numbers[1:] {
		if num > maxNum {
			maxNum = num
		}
	}
	return maxNum
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	result := make([]int, 0, len(numbers))
	existNumbers := make(map[int]struct{}, len(numbers))

	for _, num := range numbers {
		if _, ok := existNumbers[num]; !ok {
			existNumbers[num] = struct{}{}
			result = append(result, num)
		}
	}
	return result
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	n := len(slice)
	result := make([]int, n)

	for i := 0; i < n; i++ {
		result[n-1-i] = slice[i]
	}
	return result
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	result := make([]int, 0, len(numbers))

	for _, num := range numbers {
		if num%2 == 0 {
			result = append(result, num)
		}
	}
	return result
}
