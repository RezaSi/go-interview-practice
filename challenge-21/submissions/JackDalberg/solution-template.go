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
	length := len(arr)
	left, right := 0, length - 1
	if length == 0 {
	    return -1
	}
	for searchPos := length/2; right >= left; searchPos = left + (right - left)/2{
	    if arr[searchPos] == target{
	        return searchPos
	    }else if arr[searchPos] < target{
	        left = searchPos + 1
	    }else{
	        right = searchPos - 1
	    }
	}
	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	length := right - left + 1
	if length == 0 || len(arr) == 0 {
	    return -1
	}
	var searchPos int = left + (length/2)
	if target == arr[searchPos]{
	    return searchPos
	}else if target > arr[searchPos]{
	    return BinarySearchRecursive(arr, target, searchPos+1, right)
	}else{
	    return BinarySearchRecursive(arr, target, left, searchPos-1)
	}
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
    pos := BinarySearch(arr, target)
    if pos >= 0{
        return pos
    }
    length := len(arr)
    if length == 0 || target <= arr[0] {
        return 0
    }else if target >= arr[length -1]{
        return length
    }
    left, right := 0, length - 1
    for searchPos := length/2; left <= right; searchPos = left + (right - left)/2{
        if arr[searchPos] > target{
            right = searchPos - 1
        }else if arr[searchPos+1] < target{
            left = searchPos + 1
        }else{
            return searchPos + 1
        }
    }
    return -1
}
