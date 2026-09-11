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
func FindMax(nums []int) int {
    if len(nums) == 0 {
        return 0
    }
    
	maxNum := nums[0]
	for _, v := range nums[1:] {
	    if v > maxNum {
	        maxNum = v
	    }
	}
	return maxNum
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(nums []int) []int {
    res := make([]int, 0, len(nums))
	uniqMap := make(map[int]struct{}, len(nums))
	for _, v := range nums {
	    if _, exists := uniqMap[v]; !exists {
	        uniqMap[v] = struct{}{}
	        res = append(res, v)
	    }
	}
	return res
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
    newSlice := make([]int, len(slice))

	left, right := 0, len(slice)-1
	for left < right {
	    newSlice[left], newSlice[right] = slice[right], slice[left]
	    left++
	    right--
	}
	
	if left == right {
	    newSlice[left] = slice[left]
	}
	
	return newSlice
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(nums []int) []int {
	evens := make([]int, 0, len(nums))
	for _, v := range nums {
	    if v % 2 == 0 {
	        evens = append(evens, v)
	    }
	}
	return evens
}
