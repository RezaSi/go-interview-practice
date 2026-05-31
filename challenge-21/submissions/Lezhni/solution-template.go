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
    if len(arr) == 0 {
        return -1
    }

	left := 0
	right := len(arr) - 1
	
	for left <= right {
	    middle := int((left + right) / 2);
	    
	    if arr[middle] == target {
	        return middle
	    } else if target < arr[middle] {
	        right = middle - 1
	    } else {
	        left = middle + 1
	    }
	}

	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
    if len(arr) == 0 || left > right {
        return -1
    }

	middle := int((left + right) / 2)
	
	if arr[middle] == target {
	    return middle
	} else if (target < arr[middle]) {
	    return BinarySearchRecursive(arr, target, left, middle - 1)
	} else {
	    return BinarySearchRecursive(arr, target, middle + 1, right)
	}
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
    if len(arr) == 0 {
        return 0
    }

	left := 0
	right := len(arr)
	
    for left < right {
        middle := left + (right - left) / 2
        if arr[middle] < target {
            left = middle + 1
        } else {
            right = middle
        }
    }

    return left
}
