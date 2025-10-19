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

	for _, num := range numbers {
		if num > max {
			max = num
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
	uniqueNumbersMap := make(map[int]bool)
	var uniqueNumbers []int

	for _, num := range numbers {
		if _, ok := uniqueNumbersMap[num]; !ok {
			uniqueNumbersMap[num] = true
			uniqueNumbers = append(uniqueNumbers, num)
		}
	}

	return uniqueNumbers
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	// TODO: Implement this function
	sliceCopie := make([]int, len(slice))
	copy(sliceCopie, slice)

	if len(sliceCopie) == 0 {
		return []int{}
	}

	start := 0
	end := len(sliceCopie) - 1

	for start < len(sliceCopie)/2 && start < end {
		tmp := sliceCopie[start]
		sliceCopie[start] = sliceCopie[end]
		sliceCopie[end] = tmp
		start++
		end--
	}

	return sliceCopie
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	// TODO: Implement this function
	evenNumbers := []int{}
	for _, num := range numbers {
		if num%2 == 0 {
			evenNumbers = append(evenNumbers, num)
		}
	}
	return evenNumbers
}
