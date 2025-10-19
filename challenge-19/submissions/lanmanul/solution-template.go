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

	maxInt := numbers[0]
	for _, num := range numbers {
		if num > maxInt {
			maxInt = num
		}
	}
	return maxInt
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	if len(numbers) == 0 {
		return []int{}
	}
	if len(numbers) == 1 {
		return numbers
	}
	var output []int
	seen := make(map[int]struct{})

	for _, v := range numbers {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			output = append(output, v)
		}
	}

	return output
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	if len(slice) == 0 {
		return []int{}
	}
	reversed := make([]int, len(slice))
	if len(slice) == 1 {
		reversed[0] = slice[0]
		return reversed
	}

	
	for i, j := 0, len(slice)-1; i < len(slice); i, j = i+1, j-1 {
		reversed[i] = slice[j]
	}
	return reversed
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	if len(numbers) == 0 {
		return []int{}
	}
	var evenNums []int
	for _, num := range numbers {
		if num%2 == 0 {
			evenNums = append(evenNums, num)
		}
	}
	if len(evenNums) == 0 {
		return []int{}
	}
	return evenNums
}