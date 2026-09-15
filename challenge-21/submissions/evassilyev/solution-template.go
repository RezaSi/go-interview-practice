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
	la := len(arr)
	if la == 0 {
		return -1
	}
	left, right := 0, la-1

	for {
		tarr := arr[left : right+1]

		si := (len(tarr) - 1) / 2

		if target == tarr[si] {
			return left + si
		}
		if len(tarr) == 1 {
			break
		}

		if tarr[si] > target {
			if si == 0 {
				si = 1
			}
			right -= si
		} else {
			if si == 0 {
				si = 1
			}
			left += si
		}
	}
	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	if la := len(arr); la == 0 {
		return -1
	}

	tarr := arr[left : right+1]
	si := (len(tarr) - 1) / 2
	if target == tarr[si] {
		return left + si
	}
	if len(tarr) == 1 {
		return -1
	}

	if tarr[si] > target {
		if si == 0 {
			si = 1
		}
		right -= si
	} else {
		if si == 0 {
			si = 1
		}
		left += si
	}
	return BinarySearchRecursive(arr, target, left, right)
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	la := len(arr)
	if la == 0 {
		return 0
	}
	left, right := 0, la-1

	var last, lastPos int

	for {
		tarr := arr[left : right+1]

		si := (len(tarr) - 1) / 2

		if target == tarr[si] {
			return left + si
		}
		if len(tarr) == 1 {
			last = tarr[si]
			lastPos = left + si
			break
		}

		if tarr[si] > target {
			if si == 0 {
				si = 1
			}
			right -= si
		} else {
			if si == 0 {
				si = 1
			}
			left += si
		}
	}

	if target > last {
		return lastPos + 1
	}
	return lastPos
}
