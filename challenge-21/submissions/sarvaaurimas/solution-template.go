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
	// Presearch validation checks
	if len(arr) == 0 {
		return -1
	}

	// If lowest n or highest n is already out of bounds for the target
	if arr[0] > target || arr[len(arr)-1] < target {
		return -1
	}
	var (
		left  int
		right int
		mid   int
		val   int
	)
	right = len(arr) - 1
	for left <= right {
		mid = left + (right-left)/2
		val = arr[mid]
		if val == target {
			return mid
		}
		if val > target {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	// Presearch validation checks
	if len(arr) == 0 {
		return -1
	}

	// If lowest n or highest n is already out of bounds for the target
	if arr[0] > target || arr[len(arr)-1] < target {
		return -1
	}

	var (
		mid int
		val int
	)

	// If exceeded bounds return -1
	if left > right {
		return -1
	}

	mid = left + (right-left)/2
	val = arr[mid]
	// If found return it
	if val == target {
		return mid
	}

	// If more calls needed calculate new bounds and call recursively
	if val > target {
		right = mid - 1
	} else {
		left = mid + 1
	}
	return BinarySearchRecursive(arr, target, left, right)
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	// Presearch validation checks
	if len(arr) == 0 {
		return 0
	}
	if arr[0] > target {
		return 0
	}
	if arr[len(arr)-1] < target {
		return len(arr)
	}
	var (
		left  int
		right int
		mid   int
		val   int
	)
	right = len(arr) - 1
	// When they become equal its the last step
	for left <= right {
		mid = left + (right-left)/2
		val = arr[mid]
		if val == target {
			return mid
		}
		if val > target {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}

	return left
}
