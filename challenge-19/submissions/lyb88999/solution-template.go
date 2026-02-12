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
	maxValue := numbers[0]
	for i := 1; i < len(numbers); i++ {
		if numbers[i] > maxValue {
			maxValue = numbers[i]
		}
	}
	return maxValue
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	// TODO: Implement this function
	processed := make([]int, 0)
	ht := make(map[int]bool)
	for _, number := range numbers {
		if !ht[number] {
			processed = append(processed, number)
			ht[number] = true
		}
	}
	return processed
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	// TODO: Implement this function
	sliceCopy := make([]int, len(slice))
	copy(sliceCopy, slice)
	left, right := 0, len(sliceCopy)-1
	for left < right {
		sliceCopy[left], sliceCopy[right] = sliceCopy[right], sliceCopy[left]
		left++
		right--
	}
	return sliceCopy
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	filter := make([]int, 0)
	for _, number := range numbers {
		if number%2 == 0 {
			filter = append(filter, number)
		}
	}
	return filter
}
