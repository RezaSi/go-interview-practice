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
	if len(numbers) == 0{
	 return 0
	}
	max := numbers[0]
	for _ , number := range numbers[1:] {
	    if max < number{
	       max = number
	    }
	}

	return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	if len(numbers) == 0{
	 return []int{}
	}
	result := make([]int,0,len(numbers))

	result = append(result, numbers[0])
	for _, number := range numbers {
	    found := false
	    for _, inresult := range result {
	        if number == inresult {
	            found = true
	            }
	        }
	   if !found {
	       result = append(result, number)
	   }
	}
	return result
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
    if len(slice) == 0{
	 return []int{}
	}
// 	if len(slice) == 1{
// 	 return slice
// 	}
	reversed := make([]int,0,len(slice))
	for i := len(slice)-1; i >= 0 ; i-- {
	    reversed = append(reversed,slice[i])
	}
	return reversed
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	if len(numbers) == 0{
	 return []int{}
	}
	filtered := make([]int,0)
	for _, number := range numbers {
	    if number%2 == 0 {
	        filtered = append(filtered, number)
	    }
	}
	return filtered
}
