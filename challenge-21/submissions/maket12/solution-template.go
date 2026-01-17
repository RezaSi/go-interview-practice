package main

import (
	"fmt"
)

func main() {
	// Example sorted array for testing
	arr := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}

	// Test binary search
	target := 7
	index := BinarySearch(arr, target)
	fmt.Printf("BinarySearch: %d found at index %d\n", target, index)

	// Test recursive binary search
	recursiveIndex := BinarySearchRecursive(arr, target, 0, len(arr)-1)
	fmt.Printf("BinarySearchRecursive: %d found at index %d\n", target, recursiveIndex)

	// Test find insert position
	insertTarget := 8
	insertPos := FindInsertPosition(arr, insertTarget)
	fmt.Printf("FindInsertPosition: %d should be inserted at index %d\n", insertTarget, insertPos)
}

// BinarySearch performs a standard binary search to find the target in the sorted array.
// Returns the index of the target if found, or -1 if not found.
func BinarySearch(arr []int, target int) int {
	if arr == nil || len(arr) == 0 {
		return -1
	}
	if arr[0] == target {
		return 0
	}

	var left, right int = 0, len(arr) - 1
	var middle int
	for left <= right {
		middle = left + (right-left)/2
		switch {
		case arr[middle] == target:
			return middle
		case arr[middle] < target:
			left = middle + 1
		case arr[middle] > target:
			right = middle - 1
		}
	}

	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	if arr == nil || len(arr) == 0 {
		return -1
	}
	if arr[0] == target {
		return 0
	}

	if left > right {
		return -1
	}

	var middle int = left + (right-left)/2

	if arr[middle] == target {
		return middle
	}
	if arr[middle] < target {
		return BinarySearchRecursive(arr, target, middle+1, right)
	}
	if arr[middle] > target {
		return BinarySearchRecursive(arr, target, left, middle-1)
	}

	return -1
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	if arr == nil {
		return 0
	}

	left, right := 0, len(arr)-1

	for left <= right {
		middle := left + (right-left)/2
		if arr[middle] == target {
			return middle
		}
		if arr[middle] < target {
			left = middle + 1
		} else {
			right = middle - 1
		}
	}
	return left
}
