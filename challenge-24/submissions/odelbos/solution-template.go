package main

import (
	"fmt"
)

func main() {
	// Test cases
	testCases := []struct {
		nums []int
		name string
	}{
		{[]int{10, 9, 2, 5, 3, 7, 101, 18}, "Example 1"},
		{[]int{0, 1, 0, 3, 2, 3}, "Example 2"},
		{[]int{7, 7, 7, 7, 7, 7, 7}, "All same numbers"},
		{[]int{4, 10, 4, 3, 8, 9}, "Non-trivial example"},
		{[]int{}, "Empty array"},
		{[]int{5}, "Single element"},
		{[]int{5, 4, 3, 2, 1}, "Decreasing order"},
		{[]int{1, 2, 3, 4, 5}, "Increasing order"},
	}

	// Test each approach
	for _, tc := range testCases {
		fmt.Printf("Test Case: %s\n", tc.name)
		fmt.Printf("Input: %v\n", tc.nums)

		// Standard dynamic programming approach
		dpLength := DPLongestIncreasingSubsequence(tc.nums)
		fmt.Printf("DP Solution - LIS Length: %d\n", dpLength)

		// Optimized approach
		optLength := OptimizedLIS(tc.nums)
		fmt.Printf("Optimized Solution - LIS Length: %d\n", optLength)

		// Get the actual elements
		lisElements := GetLISElements(tc.nums)
		fmt.Printf("LIS Elements: %v\n", lisElements)
		fmt.Println("-----------------------------------")
	}
}

// DPLongestIncreasingSubsequence finds the length of the longest increasing subsequence
// using a standard dynamic programming approach with O(n²) time complexity.
func DPLongestIncreasingSubsequence(nums []int) int {
	n := len(nums)
	if n == 0 || n == 1 {
		return n
	}

	seq := make([]int, n)
	result := 0
	for i, val := range(nums) {
		seq[i] = 1
		for j := range(i) {
			if nums[j] < val {
				seq[i] = max(seq[i], seq[j] + 1)
			}
		}
		result = max(result, seq[i])
	}
	return result
}

// OptimizedLIS finds the length of the longest increasing subsequence
// using an optimized approach with O(n log n) time complexity.
func OptimizedLIS(nums []int) int {
	n := len(nums)
	if n == 0 || n == 1 {
		return n
	}

	return len(binarySearchLIS(nums))
}

// GetLISElements returns one possible longest increasing subsequence
// (not just the length, but the actual elements).
func GetLISElements(nums []int) []int {
	return binarySearchLIS(nums)
}

func binarySearchLIS(nums []int) []int {
	if len(nums) == 0 {
		return []int{}
	}

	seq := []int{}
	for _, num := range(nums) {
		left, right := 0, len(seq)
		for left < right {
			mid := (left + right) / 2
			if seq[mid] < num {
				left = mid + 1
			} else {
				right = mid
			}
		}

		if left == len(seq) {
			seq = append(seq, num)
		} else {
			seq[left] = num
		}
	}
	return seq
}
