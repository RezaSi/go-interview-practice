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
	left := 0
	right := len(arr) - 1 // 5

	for left <= right {
		// middle index calculation
		mid := left + (right-left)/2
		if arr[mid] == target {
			fmt.Println("target index is: ")
			return mid // target found, return index
		}

		// if value is less than teh target, increment the value
		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1 // if we get here we didn't hit above so return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
		// Use recursion to perform binary search instead
	if left > right {
		return -1 // Base case: target not found
	}

	mid := left + (right-left)/2

	if arr[mid] == target {
		return mid // Target found, return index
	} else if arr[mid] < target {
		// So instead of using a `for {...}` loop we just call the function again with the new left and right values
		return BinarySearchRecursive(arr, target, mid+1, right) // Search in the right half
	} else {
		return BinarySearchRecursive(arr, target, left, mid-1) // Search in the left half
	}

	return -1
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
		left := 0
	right := len(arr) - 1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			return mid // Target found, return index
		} else if arr[mid] < target {
			left = mid + 1 // Search in the right half
		} else {
			right = mid - 1 // Search in the left half
		}
	}
	return left // Return the insertion point
	return -1
}
