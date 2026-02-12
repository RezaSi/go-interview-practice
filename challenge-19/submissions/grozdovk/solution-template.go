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
	if len(numbers) == 0{
	    return 0
	}
	max := numbers[0]
	for _,value:= range numbers{
	    if value>max{
	        max = value
	    }
	}
	return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
	// TODO: Implement this function
	if len(numbers) == 0 {
	    return numbers
	}
	seen := make(map[int]bool)
	result:= make([]int,0, len(numbers))
	for _,value:= range numbers{
	    if !seen[value]{
	        result = append(result, value)
	        seen[value] = true
	    }
	}
	return result
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	// TODO: Implement this function
	if len(slice) == 0{
	    return slice
	}
	result:= make([]int, len(slice))
	for i,v:= range slice{
	    result[len(slice)-1-i]= v
	}
	return result
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
	if len(numbers) == 0{
	    return numbers
	}
	result:= make([]int, 0 , len(numbers))
	for _,v := range numbers{
	    if (v%2==0){
	        result = append(result,v )
	    }
	}
	return result
}
