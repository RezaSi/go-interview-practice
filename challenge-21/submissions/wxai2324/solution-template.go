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
	arrLen := len(arr)
	if arrLen == 0 {
		return -1
	}
	if arrLen == 1 {
		if target == arr[0] {
			return 0
		}
		return -1
	}
	left := 0
	right := arrLen - 1
	for left <= right {
		mid := left + (right-left)/2
		if mid >= arrLen {
			return -1
		}
		if target == arr[mid] {
			return mid
		}
		if target < arr[mid] {
			right = mid - 1
		}
		if target > arr[mid] {
			left = mid + 1
		}
	}
	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	arrLen := len(arr)

	if left > right || arrLen == 0 {
		return -1
	}
	if arrLen == 1 {
		if target == arr[0] {
			return 0
		}
		return -1
	}
	mid := left + (right-left)/2
	if left == right {
		mid = right
	}
	if target == arr[mid] {
		return mid
	}
	if target < arr[mid] {
		return BinarySearchRecursive(arr, target, left, mid-1)
	}
	return BinarySearchRecursive(arr, target, mid+1, right)
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	// 比left大 比right小 就行了 边界就是 比right还大了
	arrLen := len(arr)
	if arrLen == 0 {
		return 0
	}
	if arrLen == 1 && target == arr[0] {
		return 0
	}
	left := 0
	right := arrLen - 1
	for left <= right {
		mid := left + (right-left)/2
		if mid >= arrLen {
			return mid
		}
		if target == arr[mid] {
			return mid
		}
		if target < arr[mid] {
			right = mid - 1
		}
		if target > arr[mid] {
			left = mid + 1
		}
	}
	return right + 1
}
