package main

import (
	"fmt"
)

// This app make a binary search in sorted array by three ways:
// standard - find element
// recursive - find index of the element
// insert - find position of the element to insert,
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
	if len(arr) == 0 {
		return -1
	} else if len(arr) == 1 {
		if target == arr[0] {
			return 0
		} else {
			return -1
		}
	}
	if target < arr[0] || target > arr[len(arr)-1] {
		return -1
	}
	var mid int
	left := 0
	right := len(arr) - 1

	for left <= right {
		mid = left + (right-left)/2
		if target == arr[mid] {
			return mid
		}
		if target > arr[mid] {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	if len(arr) == 0 {
		return -1
	}
	if target < arr[0] || target > arr[len(arr)-1] {
		return -1
	}
	mid := left + (right-left)/2
	switch {
	case target == arr[mid]:
		return mid
	case left > right:
		return -1
	case target > arr[mid]:
		left = mid + 1
	case target < arr[mid]:
		right = mid - 1
	}
	return BinarySearchRecursive(arr, target, left, right)
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	left := 0
	right := len(arr) // We use len(arr) because the target could be inserted at the very end

	for left < right {
		mid := left + (right-left)/2

		if arr[mid] < target {
			// Target is to the right, move the left bound up
			left = mid + 1
		} else {
			// Target is here or to the left, move the right bound down
			right = mid
		}
	}
	return left
}
