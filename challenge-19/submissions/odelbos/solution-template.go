package main

import (
	"fmt"
	"slices"
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
    max := numbers[0]
    for _, n := range(numbers) {
        if n > max {
            max = n
        }
    }
    return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
    result := []int{}
    for _, n := range(numbers) {
        if ! slices.Contains(result, n) {
            result = append(result, n)
        }
    }
	return result
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
    size := len(slice)
    result := make([]int, size)
    for i, n := range(slice) {
        result[size - 1 - i] = n
    }
	return result
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
    result := []int{}
    for _, n := range(numbers) {
        if n % 2 == 0 {
            result = append(result, n)
        }
    }
	return result
}
