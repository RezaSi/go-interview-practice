package main

import (
	"fmt"
)

func main() {
	// Example slice for testing
	numbers := []int{30, 1, 4, 1, 5, 9, 2, 6}
	fmt.Printf("slice: %v\n", numbers)

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

	if len(numbers) <= 0 {
		return 0
	}

	max := numbers[0]
	for _, val := range numbers[1:] {
		if max < val {
			max = val
		}
	}

	// TODO: Implement this function
	return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	// TODO: Implement this function

	retSlice := make([]int, 0, len(numbers))
	retMap := make(map[int]int)

	for _, val := range numbers {
		_, ok := retMap[val]
		if !ok {
			retSlice = append(retSlice, val)
		}
		retMap[val]++
	}

	return retSlice
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	// TODO: Implement this function

	retSlice := make([]int, 0, len(slice))
	for i, _ := range slice {
		retSlice = append(retSlice, slice[len(slice)-i-1])
	}

	return retSlice
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	// TODO: Implement this function
	retSlice := make([]int, 0, len(numbers))
	for _, val := range numbers {
		if val%2 == 0 {
			retSlice = append(retSlice, val)
		}

	}

	return retSlice
}
