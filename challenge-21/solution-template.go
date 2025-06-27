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
	// TODO: Target in between test not detected -> inf loop
	l := len(arr)
	if l == 0 {
		return -1
	}
	if l == 1 {
		if arr[0] == target {
			return 0
		}
		return -1
	}

	leftEnd, rightEnd := 0, l-1
	midIdx := rightEnd / 2
	// done := false

	currLeft, currRight := leftEnd, rightEnd
	for {
		// see if we've hit the target, bingo
		if arr[midIdx] == target {
			return midIdx
		}
		// see if we've reached either end and we're done, not found
		if midIdx == leftEnd || midIdx == rightEnd {
			return -1
		}
		// go right
		if arr[midIdx] < target {
			currLeft = midIdx
			midIdx = currLeft + (currRight-currLeft)/2
			if midIdx == currLeft {
				midIdx++
			}
		} else { // go left
			currRight = midIdx
			midIdx = currLeft + (currRight-currLeft)/2
		}
	}
	/*
			Start with the middle element of the array
		If the target value equals the middle element, we're done
		If the target value is less than the middle element, search the left half
		If the target value is greater than the middle element, search the right half
		Repeat until the element is found or the search space is empty
		Binary search has a time complexity of O(log n), which is much more efficient than linear search (O(n)) for large datasets.

		Implementation Approaches
		Iterative Implementation:

		Use a loop with left and right pointers
		Calculate middle index in each iteration
		Adjust pointers based on comparison with target
		Continue until element is found or search space is exhausted

	*/
	// midIdx :=
	// if l ==  {
	// 	if arr[0] == target {
	// 		return 0
	// 	}
	// }
	return -2
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
