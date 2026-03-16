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

func BinarySearch(arr []int, target int) int {
	if len(arr) == 0 {
		return -1
	}

	low := 0
	high := len(arr) - 1

	for low <= high { //работает пока границы не схлопнулись, меняем не сам массив а индексы, которые сужжаются.
		mid := (low + high) / 2

		if arr[mid] == target {
			return mid
		}

		if arr[mid] < target { //сужжение
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}

func BinarySearchRecursive(arr []int, target int, left int, right int) int {
	if len(arr) == 0 {
		return -1
	}

	mid := (left + right) / 2

	if left > right {
		return -1
	}

	if arr[mid] == target {
		return mid
	}
	if target > arr[mid] {
		left = mid + 1
	} else {
		right = mid - 1
	}
	return BinarySearchRecursive(arr, target, left, right)
}

func FindInsertPosition(arr []int, target int) int {
	left := 0
	right := len(arr) - 1

	for left <= right {
		mid := (left + right) / 2

		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}