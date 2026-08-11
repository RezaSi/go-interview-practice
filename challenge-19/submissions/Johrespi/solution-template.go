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
	if len(numbers) == 0 {
	    return 0
	}
	
    max := numbers[0]
	
	for i := 1; i < len(numbers); i++ {
	    if numbers[i] > max {
	        max = numbers[i]
	    }
	}
	
	return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	// TODO: Implement this function
	
	if len(numbers) == 0 {
	    return []int{}
	}
	
	var res []int
	m := make(map[int]struct{})
	
	for _, e := range numbers {
	    _, exists := m[e]
	    if !exists{
	        m[e] = struct{}{}
	        res = append(res, e)
	    }
	}
	return res
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	// TODO: Implement this function
	
	if len(slice) == 0 {
	    return []int{}
	}
	
	var res []int
	
	for i := len(slice) -1; i >= 0; i--{
	    res = append(res, slice[i])
	}
	return res
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	// TODO: Implement this function
	if len(numbers) == 0 {
	    return []int{}
	}
	
	var res []int
	for _, e := range numbers {
	    if e % 2 == 0 {
	        res = append(res, e)
	    }
	}
	
	if len(res) == 0 {
	    return []int{}
	}
	
	return res
}
