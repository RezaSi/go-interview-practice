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
	// get len of arr and assert > 0
	l := len(arr)
	if l == 0 {
		return -1
	}

	// set boundaries and initial search indices
	leftEnd, rightEnd := 0, l-1
	midIdx := rightEnd / 2
	currLeft, currRight := leftEnd, rightEnd

	// loop thru mid points
	for {
		// see if we've hit the target, bingo
		if arr[midIdx] == target {
			return midIdx
		}

		// see if we've reached either end and we're done, i.e. not found
		if midIdx == leftEnd || midIdx == rightEnd {
			return -1
		}

		// check if current bounds adjacent, check both for target or done
		if currRight-currLeft == 1 {
			if arr[currLeft] == target {
				return currLeft
			}
			if arr[currRight] == target {
				return currRight
			}
			return -1
		}

		// check for left or right
		if arr[midIdx] < target {
			currLeft = midIdx // go right
		} else {
			currRight = midIdx // go left
		}

		// set new mid point
		midIdx = currLeft + (currRight-currLeft)/2
	}
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	// TODO: Implement this function
	return -1
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	// TODO: Implement this function
	return -1
}
