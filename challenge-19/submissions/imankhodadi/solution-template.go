package main

import (
	"fmt"
)

func main() {
	numbers := []int{3, 1, 4, 1, 5, 9, 2, 6}
	max := FindMax(numbers)
	fmt.Printf("Maximum value: %d\n", max)
	unique := RemoveDuplicates(numbers)
	fmt.Printf("After removing duplicates: %v\n", unique)
	reversed := ReverseSlice(numbers)
	fmt.Printf("Reversed: %v\n", reversed)
	evenOnly := FilterEven(numbers)
	fmt.Printf("Even numbers only: %v\n", evenOnly)
}
// function signature is part of the assignment
func FindMax(numbers []int) int {
	if len(numbers) == 0 {
		return 0 // default is 0
	}
	max := numbers[0]
	for _, n := range numbers[1:] {
		if n > max {
			max = n
		}
	}
	return max
}

func RemoveDuplicates(numbers []int) []int {
	if len(numbers) == 0 {
		return []int{}
	}
	seen := make(map[int]bool)
	result := make([]int, 0, len(numbers))
	for _, n := range numbers {
		if !seen[n] {
			seen[n] = true
			result = append(result, n)
		}
	}
	return result
}

func ReverseSlice(slice []int) []int {
	result := make([]int, len(slice))
	for i, v := range slice {
		result[len(slice)-1-i] = v
	}
	return result
}

func ReverseInPlace(slice []int) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func FilterEven(numbers []int) []int {
	result := make([]int, 0, len(numbers)/2)
	for _, num := range numbers {
		if num%2 == 0 {
			result = append(result, num)
		}
	}
	return result
}
