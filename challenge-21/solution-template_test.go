package main

import (
	"testing"
)

func TestBinarySearch(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{"Empty array", []int{}, 5, -1},
		{"Single element found", []int{5}, 5, 0},
		{"Single element not found", []int{5}, 7, -1},
		{"Target in the middle", []int{1, 3, 5, 7, 9}, 5, 2},
		{"Target at the beginning", []int{1, 3, 5, 7, 9}, 1, 0},
		{"Target at the end", []int{1, 3, 5, 7, 9}, 9, 4},
		{"Target not found (too small)", []int{1, 3, 5, 7, 9}, 0, -1},
		{"Target not found (too large)", []int{1, 3, 5, 7, 9}, 10, -1},
		{"Target not found (in between)", []int{1, 3, 5, 7, 9}, 4, -1},
		{"Large sorted array", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, 13, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BinarySearch(tt.arr, tt.target)
			if result != tt.expected {
				t.Errorf("BinarySearch(%v, %d) = %d, expected %d", tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

func TestBinarySearchRecursive(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{"Empty array", []int{}, 5, -1},
		{"Single element found", []int{5}, 5, 0},
		{"Single element not found", []int{5}, 7, -1},
		{"Target in the middle", []int{1, 3, 5, 7, 9}, 5, 2},
		{"Target at the beginning", []int{1, 3, 5, 7, 9}, 1, 0},
		{"Target at the end", []int{1, 3, 5, 7, 9}, 9, 4},
		{"Target not found (too small)", []int{1, 3, 5, 7, 9}, 0, -1},
		{"Target not found (too large)", []int{1, 3, 5, 7, 9}, 10, -1},
		{"Target not found (in between)", []int{1, 3, 5, 7, 9}, 4, -1},
		{"Large sorted array", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, 13, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var right int
			if len(tt.arr) > 0 {
				right = len(tt.arr) - 1
			} else {
				right = -1
			}
			result := BinarySearchRecursive(tt.arr, tt.target, 0, right)
			if result != tt.expected {
				t.Errorf("BinarySearchRecursive(%v, %d, 0, %d) = %d, expected %d",
					tt.arr, tt.target, right, result, tt.expected)
			}
		})
	}
}

func TestFindInsertPosition(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int
		target   int
		expected int
	}{
		{"Empty array", []int{}, 5, 0},
		{"Single element small", []int{5}, 3, 0},
		{"Single element bigger", []int{1}, 5, 1},
		{"Single element equal ", []int{1}, 1, 0},
		{"Insert at beginning", []int{2, 4, 6, 8}, 1, 0},
		{"Insert at end", []int{2, 4, 6, 8}, 10, 4},
		{"Insert in middle", []int{2, 4, 6, 8}, 5, 2},
		{"Target already exists (beginning)", []int{2, 4, 6, 8}, 2, 0},
		{"Target already exists (middle)", []int{2, 4, 6, 8}, 6, 2},
		{"Target already exists (end)", []int{2, 4, 6, 8}, 8, 3},
		{"Large sorted array, insert in middle", []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}, 10, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindInsertPosition(tt.arr, tt.target)
			if result != tt.expected {
				t.Errorf("FindInsertPosition(%v, %d) = %d, expected %d", tt.arr, tt.target, result, tt.expected)
			}
		})
	}
}

func TestExampleCases(t *testing.T) {
	// Example 1
	arr1 := []int{1, 3, 5, 7, 9}
	if result := BinarySearch(arr1, 5); result != 2 {
		t.Errorf("Example 1: BinarySearch(%v, 5) = %d, expected 2", arr1, result)
	}

	// Example 2
	arr2 := []int{1, 3, 5, 7, 9}
	if result := BinarySearch(arr2, 6); result != -1 {
		t.Errorf("Example 2: BinarySearch(%v, 6) = %d, expected -1", arr2, result)
	}

	// Example 3
	arr3 := []int{1, 3, 5, 7, 9}
	if result := BinarySearchRecursive(arr3, 7, 0, len(arr3)-1); result != 3 {
		t.Errorf("Example 3: BinarySearchRecursive(%v, 7, 0, %d) = %d, expected 3",
			arr3, len(arr3)-1, result)
	}

	// Example 4
	arr4 := []int{1, 3, 5, 7, 9}
	if result := FindInsertPosition(arr4, 6); result != 3 {
		t.Errorf("Example 4: FindInsertPosition(%v, 6) = %d, expected 3", arr4, result)
	}
}
