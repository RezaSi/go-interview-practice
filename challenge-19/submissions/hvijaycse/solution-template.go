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
	max_val := math.MinInt
	for _, num := range numbers {
		max_val = max(max_val, num)
	}

	return max_val
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {

	exists := map[int]bool{}
	output := make([]int, 0)

	for _, num := range numbers {
		_, exist := exists[num]
		if !exist {
			output = append(output, num)
			exists[num] = true
		}
	}
	return output
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	length := len(slice)
	reversed := make([]int, length)

	for i := range length {
		reversed[i] = slice[length-1-i]
	}
	return reversed
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	evens := make([]int, 0)

	for _, num := range numbers {

		if num%2 == 0{
			evens = append(evens, num)
		}
	}
	return evens
}
