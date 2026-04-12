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
	max := numbers[0]
	for _, val := range numbers{
	    if val > max {
	        max = val
	    }
	}
	return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	noDupes := make([]int, 0, len(numbers))
	isDupe := map[int]bool{}
	for _, val := range numbers {
	    if !isDupe[val]{
	        noDupes = append(noDupes, val)
	        isDupe[val] = true
	    }
	}
	return noDupes
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
    length := len(slice)
	reverse := make([]int, 0, length)
	for idx, _ := range slice {
	    reverse = append(reverse, slice[length - 1 - idx])
	}
	return reverse
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	filtered := make([]int, 0, len(numbers))
	for _, val := range numbers {
	    if val % 2 == 0 {
	        filtered = append(filtered, val)
	    }
	}
	return filtered
}
