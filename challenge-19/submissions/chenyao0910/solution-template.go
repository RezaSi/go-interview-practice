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
	mx := numbers[0]
	for i:=1 ; i < len(numbers);i++{
	    if(numbers[i] > mx){
	         mx = numbers[i]
	    }
	}
	return mx
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
  func RemoveDuplicates(numbers []int) []int {
      seen := make(map[int]bool)
      n := make([]int, 0, len(numbers))

      for _, num := range numbers {
          if !seen[num] {  
              seen[num] = true
              n = append(n, num)
          }
      }
      return n
  }

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
	n := make([]int, len(slice))
	for i:=len(slice) - 1 ; i >= 0 ; i--{
	    // 2 1 0
	    // 0 1 2
	    n[len(slice) - i - 1] = slice[i]
	}
	return n
}

// FilterEven returns a new slice containing only the even numbers
// from the original slice.
func FilterEven(numbers []int) []int {
    n := make([]int , 0, len(numbers))
    for i:=0 ; i < len(numbers);i++{
        if(numbers[i] % 2 == 0){
            n = append(n, numbers[i])
        }
    }
	// TODO: Implement this function
	return n
}
