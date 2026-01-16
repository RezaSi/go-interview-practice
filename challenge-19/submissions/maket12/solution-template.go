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
	if numbers == nil || len(numbers) == 0 {
	    return 0
	}
	
	max := numbers[0]
	for i := range numbers {
	    if numbers[i] > max {
	        max = numbers[i]
	    }
	}
	
	return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	newNumbers := make([]int, 0)
	numbersContain := make(map[int]bool)
	for _, val := range numbers {
	    if !numbersContain[val] {
	        newNumbers = append(newNumbers, val)
	        numbersContain[val] = true
	    }
	}
	return newNumbers
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	new := make([]int, 0, len(slice))
	for i := len(slice) - 1; i >= 0; i-- {
	    new = append(new, slice[i])
	}
	return new
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	new := make([]int, 0)
	for _, val := range numbers {
	    if val % 2 == 0 {
	        new = append(new, val)
	    }
	}
	return new
}
