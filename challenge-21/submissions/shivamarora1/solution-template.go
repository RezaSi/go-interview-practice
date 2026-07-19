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
    startIdx := 0
    endIdx := len(arr) - 1

    // Use <= to ensure we check the case where startIdx == endIdx
    for startIdx <= endIdx {
        middleElement := startIdx + (endIdx-startIdx)/2

        if arr[middleElement] == target {
            return middleElement
        } else if arr[middleElement] < target {
            // Target is in the right half, exclude middleElement
            startIdx = middleElement + 1
        } else {
            // Target is in the left half, exclude middleElement
            endIdx = middleElement - 1
        }
    }

    return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	if left > right {
		return -1
	}
	middle := left + (right-left)/2
	if arr[middle] == target {
		return middle
	} else if arr[middle] > target {
		return BinarySearchRecursive(arr, target, left, middle-1)
	} else {
		return BinarySearchRecursive(arr, target, middle+1, right)
	}
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
    left := 0
    right := len(arr) - 1
    
    // We want the smallest index where arr[idx] >= target
    insertIdx := len(arr) 

    for left <= right {
        mid := left + (right-left)/2

        if arr[mid] >= target {
            // This COULD be the insert position, but there might be 
            // an even smaller index to the left that also fits.
            insertIdx = mid
            right = mid - 1
        } else {
            // Target is definitely to the right of mid
            left = mid + 1
        }
    }

    return insertIdx
}
