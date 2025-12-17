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

	arr = []int{5}
	index = BinarySearch(arr, target)
	fmt.Printf("BinarySearch: %d not found. Return value: %d \n", target, index)

}

// BinarySearch performs a standard binary search to find the target in the sorted array.
// Returns the index of the target if found, or -1 if not found.
func BinarySearch(arr []int, target int) int {
	left := 0
	right := len(arr) - 1

	if right == -1 || target < arr[left] || target > arr[right] {
		return -1
	}

	if arr[left] == target {
		return left
	}

	if arr[right] == target {
		return right
	}

	for left <= right {
		mid := left + (right-left)/2

		// fmt.Printf("Looking between %d and %d for %d, mid: %d\n", left, right, target, mid)

		if mid == left || mid == right {
			return -1
		}

		switch {
		case target < arr[mid]:
			right = mid
		case target > arr[mid]:
			left = mid
		default:
			return mid
		}

	}

	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {

	if right == -1 || left > right || target < arr[left] || target > arr[right] {
		return -1
	}

	if arr[left] == target {
		return left
	}

	if arr[right] == target {
		return right
	}

	mid := left + (right-left)/2

	if mid == left || mid == right {
		return -1
	}

	switch {
	case target < arr[mid]:
		right = mid
	case target > arr[mid]:
		left = mid
	default:
		return mid
	}

	return BinarySearchRecursive(arr, target, left, right)
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	left := 0
	right := len(arr) - 1

	if right == -1 {
		return 0
	}

	if target < arr[left] {
		return left
	}

	if target > arr[right] {
		return right + 1
	}

	for left <= right {
		mid := left + (right-left)/2

		if mid == left || mid == right {
			if arr[mid] < target {
				return left + 1
			} else {
				return right - 1
			}
		}

		switch {
		case target < arr[mid]:
			right = mid
		case target > arr[mid]:
			left = mid
		default:
			return mid
		}

	}

	return 0

}
