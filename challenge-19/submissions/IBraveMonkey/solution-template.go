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
	
	for _, num := range numbers {
	    if num > max {
	        max=num
	    }
	}
	
	return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	hMap := make(map[int]struct{})
	res := make([]int, 0, len(numbers))
	
	for _, num := range numbers {
	    if _, exists := hMap[num]; !exists {
	        hMap[num]=struct{}{}
	        res=append(res, num)
	    }
	}

	return res
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	newSlice := make([]int, 0, len(slice))

	for i := len(slice) - 1; i >= 0; i-- {
		newSlice = append(newSlice, slice[i])
	}

	return newSlice
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	newSlice := make([]int, 0, len(numbers))
	
	for _, num := range numbers {
	    if num % 2 == 0 {
	        newSlice = append(newSlice, num)
	    }    
	}
	
	return newSlice
}
