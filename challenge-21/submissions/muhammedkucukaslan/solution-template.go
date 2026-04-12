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

func BinarySearch(arr []int, target int) int {
	if len(arr) == 0 {
		return -1
	}

	minIndex := 0
	maxIndex := len(arr) - 1

	for minIndex <= maxIndex {
		midIndex := (maxIndex + minIndex) / 2

		if arr[midIndex] == target {
			return midIndex
		}

		if arr[midIndex] > target {
			maxIndex = midIndex - 1
		} else {
			minIndex = midIndex + 1
		}
	}

	return -1
}

func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	if len(arr) == 0 {
		return -1
	}

	if left > right {
		return -1
	}

	midIndex := (left + right) / 2

	if arr[midIndex] == target {
		return midIndex
	}

	if arr[midIndex] > target {
		return BinarySearchRecursive(arr, target, left, midIndex-1)

	}
	if arr[midIndex] < target {
		return BinarySearchRecursive(arr, target, midIndex+1, right)
	}

	return -1
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	left, right := 0, len(arr)

	for left < right {
		mid := (left + right) / 2

		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid
		}
	}

	return left
}

