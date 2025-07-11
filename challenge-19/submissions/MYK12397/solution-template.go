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
	if len(numbers) > 0 {
	    max_val := -1
	    
	    for _,val := range numbers {
	        if val > max_val {
	            max_val = val
	        }
	    }
	    return max_val
	}
	return 0
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	// TODO: Implement this function
	if len(numbers) <= 1 {
	    return numbers
	}
	checkDuplicate := make(map[int]bool)
	idx := 0
	for _,val := range numbers {
	    if !checkDuplicate[val] {
	        checkDuplicate[val] = true
	        numbers[idx] = val
	        idx++
	    }
	}
	return numbers[:idx]
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	// TODO: Implement this function
	arr := make([]int,len(slice))
	copy(arr,slice)
	if len(arr) <= 1 {
	    return arr
	}
    for i, j := 0, len(arr)-1; i <= j; i, j = i+1, j-1 {
        arr[i], arr[j] = arr[j], arr[i]
    }
    
	return arr
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	// TODO: Implement this function
	slice := make([]int,0)
	
	for _,val := range numbers {
	    if val%2 == 0 {
	        slice = append(slice,val)
	    }
	}
	return slice
}
