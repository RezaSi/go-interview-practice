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
	for _, v := range numbers[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	valSet := map[int]struct{}{}
	result := []int{}
	for _, n := range numbers {
		if _, ok := valSet[n]; ok {
			continue
		}
		valSet[n] = struct{}{}
		result = append(result, n)
	}

	return result
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	lenth := len(slice)
	if lenth == 0 {
		return slice
	}
	reversed := make([]int, lenth, lenth)
	copy(reversed, slice)
	for i := 0; i < lenth/2; i++ {
		reversed[i], reversed[lenth-i-1] = slice[lenth-i-1], slice[i]
	}

	return reversed
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	even := []int{}
	for _, n := range numbers {
		if n%2 == 0 {
			even = append(even, n)
		}
	}
	return even
}
