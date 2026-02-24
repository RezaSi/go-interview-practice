package main

import (
	"fmt"
)

func main() {
	arr := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}

	target := 7
	index := BinarySearch(arr, target)
	fmt.Printf("BinarySearch: %d found at index %d\n", target, index)

	recursiveIndex := BinarySearchRecursive(arr, target, 0, len(arr)-1)
	fmt.Printf("BinarySearchRecursive: %d found at index %d\n", target, recursiveIndex)

	for i := 0; i <= 20; i += 2 {
		insertPos := FindInsertPosition(arr, i)
		fmt.Printf("FindInsertPosition: %d should be inserted at index %d\n", i, insertPos)
	}
}
func BinarySearch(arr []int, target int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		m := left + (right-left)/2
		if arr[m] == target {
			return m
		} else if target < arr[m] {
			right = m - 1
		} else {
			left = m + 1
		}
	}
	return -1
}

func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	if left > right {
		return -1
	}
	m := left + (right-left)/2
	if arr[m] == target {
		return m
	} else if target < arr[m] {
		return BinarySearchRecursive(arr, target, left, m-1)
	} else {
		return BinarySearchRecursive(arr, target, m+1, right)
	}
}

func FindInsertPosition(arr []int, target int) int {
	left, right := 0, len(arr)
	for left < right {
		mid := left + (right-left)/2
		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return left
}