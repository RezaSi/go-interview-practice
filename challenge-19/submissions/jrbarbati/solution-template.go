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
    
    max := math.MinInt
    
    for _, val := range numbers {
        if val > max {
            max = val
        }
    }
    
	return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
    seen := make(map[int]struct{})
    unique := make([]int, 0)
    
    for _, val := range numbers {
        if _, ok := seen[val]; !ok {
            unique = append(unique, val)
        }
        
        seen[val] = struct{}{}
    }
    
	return unique
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
    var j int
    reversed := make([]int, len(slice))

    for i := len(slice)-1; i >= 0; i-- {
        reversed[j] = slice[i]
        j++
    }    
    
	return reversed
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
    evens := make([]int, 0)
    
    for _, val := range numbers {
        if val%2 == 0 {
            evens = append(evens, val)
        }
    }
    
	return evens
}
