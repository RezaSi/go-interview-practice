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
	low := 0
	high := len(arr)-1
	
	for low <= high {
	    mid := (low + high) / 2
	    
	    shoot := arr[mid]
	    
	    if shoot == target {
	        return mid
	    }
	    
	    if shoot > target {
	        high = mid - 1
	    }
	    
	    if shoot < target {
	        low = mid + 1
	    }
	}
	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {

    // 1. Base case: interval is invalid.
    if left > right {
        return -1
    }

    // 2. Calculate middle index.
    mid := left + (right-left)/2

    // 3. Check if target was found.
    if arr[mid] == target {
        return mid
    }

    // 4. Search the left half.
    if target < arr[mid] {
        return BinarySearchRecursive(arr, target, left, mid-1)
    }

    // 5. Search the right half.
    return BinarySearchRecursive(arr, target, mid+1, right)
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
    for i, item := range arr {
        if target <= item {
            return i
        }
    }

    return len(arr)
}
