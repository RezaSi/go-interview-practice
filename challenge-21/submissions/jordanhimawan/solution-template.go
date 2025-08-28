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
    left, right := 0, len(arr)-1
	for left <= right {
    	mid := left + (right-left)/2
    	if arr[mid] == target {
	        return mid
	    }
	    if arr[mid] > target {
	        right = mid - 1
	    } else {
	        left = mid + 1
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
	// TODO: Implement this function
	mid := left + (right-left)/2
	index := -1
	if arr[mid] > target {
	    index = BinarySearchRecursive(arr, target, left, mid - 1)
	} else if arr[mid] < target {
	    index = BinarySearchRecursive(arr, target, mid + 1, right)
	} else if arr[mid] == target {
	    index = mid
	}
	return index
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	// TODO: Implement this function
	if len(arr) == 0{
	    return 0
	}
    low, high := 0, len(arr)

    for low < high {
        mid := (low + high) / 2
        if arr[mid] < target {
            low = mid + 1
        } else {
            high = mid
        }
    }
    return low
}
