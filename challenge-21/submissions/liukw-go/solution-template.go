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
	n, ret := len(arr), -1
	l, r := 0, n - 1
	for l <= r {
	    m := (l + r) >> 1
	    if arr[m] >= target {
	        r = m - 1
	        ret = m
	    } else {
	        l = m + 1
	    }
	}
	if ret == -1 { return -1 }
	if target != arr[ret] { return -1 }
	return ret
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	// TODO: Implement this function
	l, r, ret := left, right, -1
	for l <= r {
	    m := (l + r) >> 1
	    if arr[m] >= target {
	        r = m - 1
	        ret = m
	    } else {
	        l = m + 1
	    }
	}
	if ret == -1 { return -1 }
	if target != arr[ret] { return -1 }
	return ret
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	// TODO: Implement this function
	n, ret := len(arr), len(arr)
	l, r := 0, n - 1
	for l <= r {
	    m := (l + r) >> 1
	    if arr[m] >= target {
	        r = m - 1
	        ret = m
	    } else {
	        l = m + 1
	    }
	}
	return ret
}
