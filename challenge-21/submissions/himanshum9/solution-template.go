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
	l := 0
	r := len(arr) - 1
	for l <= r {
	    mid := (l + r) / 2
	    if target < arr[mid] {
	        r = mid - 1
	    } else if target > arr[mid] {
	        l = mid + 1
	    } else {
	        return mid
	    }
	}
	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	// TODO: Implement this function
	l := left
	r := right
	if l > r {
	    return -1
	}
	
	mid := (l + r) / 2
	if arr[mid] == target {
	    return mid
	}
	if target < arr[mid] {
	    r = mid - 1
	    return BinarySearchRecursive(arr,target,l,r)
	}
	
	if target > arr[mid] {
	    l = mid + 1
	    return BinarySearchRecursive(arr,target,l,r)
	}
	return -1
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	// TODO: Implement this function
	for i,j := range arr {
	    if target <= j {
	        return i
	    }
	}
	return len(arr)
}
