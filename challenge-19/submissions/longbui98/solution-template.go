package main

import (
	"fmt"
	"sort"
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
	sort.Ints(numbers)
	return numbers[len(numbers)-1]
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	// TODO: Implement this function
	valueDupMap := make(map[int]int)
	result := make([]int, 0, len(numbers))
	for _, num := range numbers {
		if _, ok := valueDupMap[num]; !ok {
			valueDupMap[num] = 1
			result = append(result, num)
		}
	}
	return result
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	// TODO: Implement this function
	reverse := make([]int, 0, len(slice))
	for i := len(slice) - 1; i >= 0; i-- {
		reverse = append(reverse, slice[i])
	}
	return reverse
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	// TODO: Implement this function
	even := make([]int, 0, len(numbers))
	for _, n := range numbers {
		if n%2 == 0 {
			even = append(even, n)
		}
	}
	return even
}
