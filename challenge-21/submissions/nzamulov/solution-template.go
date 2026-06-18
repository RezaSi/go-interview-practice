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
	l, r := -1, len(arr)
	for l+1 < r {
	    m := (l + r) >> 1
	    if arr[m] < target {
	        l = m
	    } else {
	        r = m
	    }
	}
	if r < len(arr) && arr[r] == target {
	    return r
	}
	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
    if left > right {
        return -1
    }
	m := (left + right) >> 1
	if arr[m] == target {
	    return m
	}
	if left == right {
	    return -1
	}
	if arr[m] < target {
	    return BinarySearchRecursive(arr, target, m+1, right)
	}
	return BinarySearchRecursive(arr, target, left, m)
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	l, r := -1, len(arr)
	for l+1 < r {
	    m := (l + r) >> 1
	    if arr[m] < target {
	        l = m
	    } else {
	        r = m
	    }
	}
	return r
}
