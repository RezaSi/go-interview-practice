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
    
    mx := numbers[0]
    for _, n := range numbers {
        mx = max(mx, n)
    }
	return mx
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
    new := make([]int, 0, len(numbers))
    added := make(map[int]bool, len(numbers))
    for _, n := range numbers {
        if added[n] {
            continue
        }
        added[n] = true
        new = append(new, n)
    }
	return new
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
    new := make([]int, 0, len(slice))
    for i := len(slice)-1; i >= 0; i-- {
        new = append(new, slice[i])
    }
	return new
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	new := make([]int, 0, len(numbers))
    for _, n := range numbers {
        if n % 2 != 0 {
            continue
        }
        new = append(new, n)
    }
	return new
}
