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
	if len(arr) == 0 || target < arr[0] || arr[len(arr)-1] < target {
		return -1
	}

	L := 0
	R := len(arr)
	var M int

	for R-L > 1 {
		M = (L + R) / 2
		if arr[M] <= target {
			L = M
		} else {
			R = M
		}
	}

	if arr[L] == target {
		return L
	} else {
		return -1
	}
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	if right < left || target < arr[left] || arr[right] < target {
		return -1
	}
	if arr[left] == target {
		return left
	} else if arr[right] == target {
		return right
	} else if right-left < 2 {
		return -1
	}

	mid := (left + right) / 2

	if arr[mid] <= target {
		return BinarySearchRecursive(arr, target, mid, right)
	} else {
		return BinarySearchRecursive(arr, target, left, mid-1)
	}
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	if len(arr) == 0 || target < arr[0] {
		return 0
	} else if arr[len(arr)-1] < target {
		return len(arr)
	}

	L := 0
	R := len(arr)
	var M int

	for R-L > 1 {
		M = (L + R) / 2
		if arr[M] <= target {
			L = M
		} else {
			R = M
		}
	}

	if arr[L] == target {
		return L
	} else {
		return R
	}
}
