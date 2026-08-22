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
    l := len(arr)
    if l == 0 { return -1 }
    
    lower,upper := 0, l-1
    if lower == upper {
        if arr[0] == target { return 0 } else { return -1 }
    }
    
    if arr[lower] > target || arr[upper] < target {
        return -1
    }
    
    result := -1
    for lower,upper := 0,l-1; lower != upper; {
	    
	    switch {
            case arr[lower] == target:
                return lower
            case arr[upper] == target:
                return upper
        }
        
	    diff := upper - lower
	    if diff == 1 {
	        return -1
	    }
	    
        halfway := diff / 2
        if arr[lower + halfway] == target {
            return lower + halfway
        }
        
	    switch {
	        case arr[lower + halfway] > target:
	            upper = lower + halfway
            case arr[lower + halfway] < target:
                lower = lower + halfway
            default:
                break
	    }
	    
	}
	
	return result
}

// BinarySearchRecursive performs binary search using recursion.
// Returns the index of the target if found, or -1 if not found.
func BinarySearchRecursive(arr []int, target int, left int, right int) int {
    l := len(arr)
    switch {
        case l == 0:
            return -1
        case left < 0:
            return -1
        case right > l + 1:
            return -1
        case left > right:
            return -1
        case arr[right] < target:
            return -1
        case arr[left] > target: return -1
    }
    
    diff := right - left
    var halfway int
	switch {
	    case arr[left] == target: return left
        case arr[right] == target: return right
        case diff == 0 || diff == 1: return -1
        default:
            halfway = left + diff/2
	}

    if arr[halfway] > target {
        return BinarySearchRecursive(arr, target, left, halfway)
    } else {
        return BinarySearchRecursive(arr, target, halfway, right)
    }
}

// FindInsertPosition returns the index where the target should be inserted
// to maintain the sorted order of the array.
func FindInsertPosition(arr []int, target int) int {
	l := len(arr)
    if l == 0 { return 0 }
    
    for lower,upper := 0,l-1; lower != upper; {
	    diff := upper - lower
	    switch {
	        case arr[lower] == target: return lower
            case arr[lower] > target: return lower
            case arr[upper] == target: return upper
            case arr[upper] < target: return upper + 1
            case diff == 1 && 
                    arr[lower] < target && 
                    arr[upper] > target:
                return lower + 1
	    }
	    
        halfway := lower + diff/2
        
	    switch {
	        case arr[halfway] == target:
	            return halfway
	        case arr[halfway] > target:
	            upper = halfway
            case arr[halfway] < target:
                lower = halfway
	    }
	    
	}
	
	return 0
}
