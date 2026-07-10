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
	if (len(numbers) == 0) {
	    return 0
	}
	res := numbers[0]
	for _, val := range numbers {
	    res = max(res, val)
	}
	return res
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	m := make(map[int]bool)
	res := make([]int, 0, len(numbers))
	for _, val := range numbers {
	    if _, ok := m[val]; !ok {
	        m[val] = true
	        res = append(res, val)
	    }
	}
	return res
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	res := make([]int, 0, len(slice))
	for i := range len(slice) {
	    res = append(res, slice[len(slice)-1-i])
	}
	return res
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	res := make([]int, 0, len(numbers))
	for _, val := range numbers {
	    if val % 2 == 0 {
	        res = append(res, val)
	    }
	}
	return res
}
