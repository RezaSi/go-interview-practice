package main

import (
	"fmt"
	"math"
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
    
	ret := math.MinInt
	for _, n := range numbers {
	    if ret < n {
	        ret = n
	    }
	}
	
	return ret
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
    ret := make([]int, 0, len(numbers))
	appeared := map[int]struct{}{}
	for _, n := range numbers {
	    if _, ok := appeared[n]; !ok {
	        ret = append(ret, n)
	        appeared[n] = struct{}{}
	    }
	}
	return ret
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	ret := make([]int, 0, len(slice))
	for i := range slice {
	    ret = append(ret, slice[len(slice)-i-1])
	}
	return ret
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
    ret := make([]int, 0, len(numbers))
	for _, n := range numbers {
	    if n % 2 == 0 {
	        ret = append(ret, n)
	    }
	}
	return ret
}
