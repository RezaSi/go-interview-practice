package main

import (
	"fmt"
)

/*
*
https://app.gointerview.dev/challenge/21
*/
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
	if len(arr) == 1 {
		if arr[0] == target {
			return 0
		}
		return -1
	}

	l := 0
	r := len(arr) - 1
	var mid int
	for l <= r {
		mid = (l + r) / 2
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	// TODO: Implement this function
	if len(arr) == 0 || target < arr[left] || target > arr[right] {
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
  if target == arr[r] {
    return r
  }
	if target > arr[r] {
		return len(arr)
	}
	if target < arr[l] {
		return 0
	}
	var mid int
	for l <= r {
		tmp := r - l
		mid = (l + r) / 2
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			if tmp == 1 {
				return r
			}
			l = mid
		} else {
			if tmp == 1 {
				return r
			}
			r = mid
		}
	}
	return -1
}