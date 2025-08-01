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
	// TODO: Implement this function
	maxNumber := 0
	if len(numbers) > 0 {
		maxNumber = numbers[0]
	}

	for i := range numbers {
		if numbers[i] > maxNumber {
			maxNumber = numbers[i]
		}
	}

	return maxNumber
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	// TODO: Implement this function
	uniqueMap := make(map[int]bool, len(numbers))
	uniqueSlice := []int{}
	for i := 0; i < len(numbers); i++ {
		if _, ok := uniqueMap[numbers[i]]; ok {
			continue
		} else {
			uniqueMap[numbers[i]] = true
			uniqueSlice = append(uniqueSlice, numbers[i])
		}
	}

	return uniqueSlice
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	// TODO: Implement this function
	reverseSlice := make([]int, len(slice))

	for i := len(slice) - 1; i >= 0; i-- {
		reverseSlice[len(slice)-1-i] = slice[i]
	}

	return reverseSlice
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	// TODO: Implement this function
	evenSlice := []int{}

	for i := range numbers {
		if numbers[i]%2 == 0 {
			evenSlice = append(evenSlice, numbers[i])
		}
	}

	return evenSlice
}
