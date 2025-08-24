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

	maxEl := numbers[0]
	for _, num := range numbers[1:] {
		if num > maxEl {
			maxEl = num
		}
	}

	return maxEl

}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	if len(numbers) == 0 {
		return []int{}
	}

	m := make(map[int]bool, len(numbers))
	uniqueNumbers := []int{}

	for _, number := range numbers {
		if _, val := m[number]; !val {
			m[number] = true
			uniqueNumbers = append(uniqueNumbers, number)
		}
	}

	return uniqueNumbers
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	if len(slice) == 0 {
		return []int{}
	}
	reverseSlice := []int{}

	for i := len(slice) - 1; i >= 0; i-- {
		reverseSlice = append(reverseSlice, slice[i])
	}

	return reverseSlice
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	if len(numbers) == 0 {
		return []int{}
	}

	evenNumbers := []int{}

	for _, number := range numbers {
		if number%2 == 0 {
			evenNumbers = append(evenNumbers, number)
		}
	}

	return evenNumbers
}
