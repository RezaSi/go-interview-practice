package main

import (
	"fmt"
)

func main() {
	// Example sorted array for testing
	arr := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 28}

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
	l, r := 0, len(arr)-1

	for l <= r {
		m := l + (r-l)/2
		fmt.Println(m)
		if arr[m] == target {
			return m
		} else if arr[m] < target {
			l = m + 1
		} else {
			r = m - 1
		}
	}
	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, l int, r int) int {
	if l > r {
		return -1
	}

	m := l + (r-l)/2

	if arr[m] == target {
		return m
	} else if arr[m] < target {
		l = m + 1
	} else {
		r = m - 1
	}

	return BinarySearchRecursive(arr, target, l, r)
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	l, r := 0, len(arr)-1

	for l <= r {
		m := l + (r-l)/2

		if arr[m] == target {
			return m
		} else if arr[m] < target {
			l = m + 1
		} else {
			r = m - 1
		}
	}
	return l;
}
