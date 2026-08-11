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
	
	length := len(arr)
	
	if length == 0 {
	    return -1
	}
	
	start := 0
	end := length - 1
	mid := end / 2

	if target == arr[start] {
	    return start
	}
	
	if target == arr[end] {
	    return end
	}
		
	for start != mid  {
	    
	    if target == arr[mid] {
	        return mid
	    }
	    
	    if target < arr[mid] {
	        end = mid
	        mid /= 2
	    } else {
	        start = mid
	        mid = mid + ((end - mid) / 2)
	    }
	}
	
	return -1
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	// TODO: Implement this function
	
	length := len(arr)
	if length == 0 {
	    return -1
	}
	
	if target == arr[left] {
	    return left
	}
	
	if target == arr[right] {
	    return right
	}
	
	mid := left + ((right - left) / 2)
	
	if target == arr[mid] {
	    return mid
	}
	
	if left == mid {
	    return -1
	}
	
	if(target < arr[mid]) {
	    return BinarySearchRecursive(arr, target, left, mid)
	} else {
	    return BinarySearchRecursive(arr, target, mid, right)
	}
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	// TODO: Implement this function
	
	length := len(arr)
	start := 0
	end := length - 1
	mid := end / 2
	
	if length == 0 {
	    return 0
	}
	
	if target > arr[end] {
	    return end + 1
	}
	
	if target == arr[start] {
	    return start
	}
	
	if target == arr[end] {
	    return end
	}
	
	for start != mid  {

	    if target == arr[mid] {
	        return mid
	    }
	    
	    if target < arr[mid] {
	        end = mid
	        mid /= 2
	    } else {
	        start = mid
	        mid = mid + ((end - mid) / 2)
	    }
	}
	
	if mid == 0 {
	    return 0
	} else {
	   if(target < arr[mid]) {
	       return mid - 1
	   } else {
	       return mid + 1
	   }
	}
}
