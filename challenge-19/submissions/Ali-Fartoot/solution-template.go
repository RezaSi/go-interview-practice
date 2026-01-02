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
            max = num
        }
    }
    return max
}

// RemoveDuplicates returns a new slice with duplicate values removed,
// preserving the original order of elements.
func RemoveDuplicates(numbers []int) []int {
    if len(numbers) == 0 {
        return []int{}
    }
    
    // Use a map to track seen values
    seen := make(map[int]bool)
    result := make([]int, 0)
    
    for _, num := range numbers {
        if !seen[num] {
            seen[num] = true
            result = append(result, num)
        }
    }
    return result
}

// ReverseSlice returns a new slice with elements in reverse order.
func ReverseSlice(slice []int) []int {
    if len(slice) == 0 {
        return []int{}
    }
    
    // Create a new slice to avoid modifying the original
    reversed := make([]int, len(slice))
    
    // Copy elements in reverse order
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
    
    even := make([]int, 0)
    
    for _, num := range numbers {
        if num%2 == 0 {
            even = append(even, num)
        }
    }
    return even
}