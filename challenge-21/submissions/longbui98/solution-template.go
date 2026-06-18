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
	lenArr := len(arr) - 1
	start, end := 0, lenArr
	for start <= end {
		mid := (start + end) / 2
		if arr[mid] == target {
			return mid
		} else if arr[mid] > target {
			end = end - 1
		} else {
			start = start + 1
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
	if arr[mid] == target {
		return mid
	} else if arr[mid] > target {
		return BinarySearchRecursive(arr, target, left, right-1)
	}
	return BinarySearchRecursive(arr, target, left+1, right)

}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	// TODO: Implement this function
	for i := 0; i < len(arr); i++ {
		if (i == 0 && target <= arr[i]) || target == arr[i] {
			return i
		}
		if i == len(arr)-1 && target > arr[i] {
			return i + 1
		}
		if i < len(arr)-1 {
			if target > arr[i] && target < arr[i+1] {
				return i + 1
			}
		}
	}
	return 0
}
