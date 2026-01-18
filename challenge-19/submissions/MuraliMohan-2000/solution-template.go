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
	if len(numbers) > 0 {
		return slices.Max(numbers)
	}
	return 0
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	if len(numbers) < 2 {
		return numbers
	}

	hashSet := make(map[int]struct{})
	var res []int

	for _, i := range numbers {
		if _, ok := hashSet[i]; !ok {

			hashSet[i] = struct{}{}
			res = append(res, i)

		}
	}

	return res
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	n := len(slice)

	if n < 2 {
		return append([]int{}, slice...)
	}

	res := make([]int, n)

	for i := 0; i < n; i++ {
		res[i] = slice[n-1-i]
	}
	return res

}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	return slices.DeleteFunc(numbers, func(n int) bool {
		return n%2 != 0
	})
}
