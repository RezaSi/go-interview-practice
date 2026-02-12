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
	// TODO: Implement this function
	if len(arr) == 0 || target > arr[len(arr)-1] || target < arr[0] {
		return -1
	}

	l := 0
	r := len(arr) - 1

	for l <= r {
		mid := l + (r-l)/2

		if target == arr[mid] {
			return mid
		}
		if target < arr[mid] {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	// TODO: Implement this function
	if left > right {
		return -1
	}

	mid := (left + right) / 2

	if target == arr[mid] {
		return mid
	} else if target < arr[mid] {
		return BinarySearchRecursive(arr, target, left, mid-1)
	} else {
		return BinarySearchRecursive(arr, target, mid+1, right)
	}
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	// TODO: Implement this function
	if len(arr) == 0 {
		return 0
	}

	l := 0
	r := len(arr) - 1

	for l <= r {
		mid := l + (r-l)/2

		if target == arr[mid] {
			return mid
		} else if target < arr[mid] {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return l
}
