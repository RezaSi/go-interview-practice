package main

import (
	"fmt"
	"sort"
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
	left := 0
	right := len(arr) - 1
	var mid int

	for right >= left {

		mid = left + ((right - left) / 2)

		switch {
		
		case arr[mid] == target:
			return mid
		
		case arr[mid] > target:
			right = mid - 1

		case arr[mid] < target:
			left = mid + 1
		}
	}
	return -1

}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
    
	var mid int

	for right >= left {

		mid = left + ((right - left) / 2)

		switch {
		
		case arr[mid] == target:
			return mid
		
		case arr[mid] > target:
			return BinarySearchRecursive(arr, target, left, mid - 1)

		case arr[mid] < target:
			return BinarySearchRecursive(arr, target, mid + 1, right)
		}
	}
	return -1
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	
	i := sort.Search(len(arr), func(i int) bool { return arr[i] >= target })
	if i < len(arr) && arr[i] == target {
		return i
	} else {
		return i
	}
	
	return -1
}
