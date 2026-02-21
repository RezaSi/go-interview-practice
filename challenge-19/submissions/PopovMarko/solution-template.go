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
	// Check for emty slice
	if len(numbers) == 0 {
		return 0
	}

	// Variable for store the max value
	max := numbers[0]

	// Loop through the slice elemenes with max
	for _, n := range numbers {
		if n > max {
			max = n
		}
	}
	return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	// Check for emty slice
	if len(numbers) == 0 {
		return []int{}
	}
	//New slice for store unic values
	res := []int{}

	// Map [int]bool to store wether elements seen in slice
	duplicate := make(map[int]bool)

	// Loop throght the numbers, check for duplicates and store in new slice
	for _, n := range numbers {
		if !duplicate[n] {
			duplicate[n] = true
			res = append(res, n)
		}
	}
	return res
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	// Check for emty slice
	if len(slice) == 0 {
		return []int{}
	}

	// New slice for store revesed elemenst
	res := make([]int, 0, len(slice))

	// Check for one element slice
	if len(slice) == 1 {
		res = append(res, slice[0])
		return res
	}

	// Loop the slice in reverce order and write elemenes to new slice
	for i := len(slice) - 1; i >= 0; i-- {
		res = append(res, slice[i])
	}
	return res
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	// Check for emty slice
	if len(numbers) == 0 {
		return []int{}
	}

	// New slice to store even elemenes
	res := make([]int, 0)

	// Loop numbers and check for even
	for _, n := range numbers {
		if n%2 == 0 {
			res = append(res, n)
		}
	}
	return res
}
