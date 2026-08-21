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
	n := len(numbers)

	if n == 0 {
		return 0
	}

	maxNum := numbers[0]
	for _, v := range numbers {
		maxNum = max(maxNum, v)
	}

	return maxNum
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	ans := make([]int, 0, len(numbers))

	m := make(map[int]int)

	for _, v := range numbers {
		m[v]++

		if m[v] == 1 {
			ans = append(ans, v)
		}
	}
	return ans
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {

	var ans []int = make([]int, 0, len(slice))

	for i := len(slice) - 1; i >= 0; i-- {
		ans = append(ans, slice[i])
	}

	return ans
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {

	ans := make([]int, 0, len(numbers))

	for _, v := range numbers {
		if v%2 == 0 {
			ans = append(ans, v)
		}
	}
	return ans

}
