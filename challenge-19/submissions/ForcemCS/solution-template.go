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
	for _, num := range numbers[1:] {
		if num > max {
			max = num
		}
	}

	return max
}

// 删除切片中的重复项，但是保留原来的顺序
func RemoveDuplicates(numbers []int) []int {

	if len(numbers) == 0 {
		return []int{}
	}

	seen := make(map[int]bool)
	var result []int

	for _, num := range numbers {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}

	return result
}

// 创建并返回一个逆序的新切片
func ReverseSlice(slice []int) []int {
	if len(slice) == 0 {
		return []int{}
	}
	result := make([]int, len(slice))

	for i, s := range slice {
		result[len(slice)-1-i] = s
	}

	return result
}

func FilterEven(numbers []int) []int {
	// TODO: Implement this function
	if len(numbers) == 0 {
	    return []int{}
	}
	
	var res []int
	for _, e := range numbers {
	    if e % 2 == 0 {
	        res = append(res, e)
	    }
	}
	
	if len(res) == 0 {
	    return []int{}
	}
	
	return res
}
