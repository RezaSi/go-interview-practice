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
leftIndex := 0
	rightIndex := len(arr)
	var midIndex int

	if rightIndex <= 0 {
		return -1
	}

	if rightIndex == 1 {
		if arr[0] == target {
			return 0
		} else {
			return -1
		}
	}

	for leftIndex <= rightIndex {
		midIndex = leftIndex + (rightIndex-leftIndex)/2
		if midIndex >= len(arr) {
			return -1
		}

		if arr[midIndex] == target {
			return midIndex
		} else if arr[midIndex] < target {
			leftIndex = midIndex + 1
		} else {
			rightIndex = midIndex - 1
		}
	}

	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
var midIndex int

	if right < 0 {
		return -1
	}

	if right == 0 {
		if arr[0] == target {
			return 0
		} else {
			return -1
		}
	}

	midIndex = left + (right-left)/2

	if left > right {
		return -1
	}

	if arr[midIndex] == target {
		return midIndex
	} else if arr[midIndex] < target {
		left = midIndex + 1
		return BinarySearchRecursive(arr, target, left, right)
	} else {
		right = midIndex - 1
		return BinarySearchRecursive(arr, target, left, right)
	}

}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
leftIndex := 1
	rightIndex := len(arr) - 1
	var midIndex int

	if rightIndex < 0 {
		return 0
	}

	if rightIndex == 0 {
		if arr[0] == target {
			return 0
		} else {
			return 1
		}
	}

	for leftIndex <= rightIndex {
		midIndex = leftIndex + (rightIndex-leftIndex)/2

		if arr[midIndex] == target {
			return midIndex
		} else if arr[midIndex] < target {
			leftIndex = midIndex + 1
		} else {
			rightIndex = midIndex - 1
		}
	}

if leftIndex == 1 {
    return 0
}

	return leftIndex
}
