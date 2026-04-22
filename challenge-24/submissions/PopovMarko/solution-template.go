package main

import (
	"fmt"
	"slices"
	"sort"
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
	// Check for empty slice in parameters
	length := len(nums)
	if length == 0 {
		return 0
	}

	// Initialization of the helper slice dp
	dp := make([]int, length)
	for i := range length {
		dp[i] = 1
	}

	// Fwo for cycles to calculate helper slice with lengthes of possible sequenses
	for i := 1; i < length; i++ {
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] {
				tmp := dp[j] + 1
				if tmp > dp[i] {
					dp[i] = tmp
				}
			}
		}
	}

	return slices.Max(dp)
}

// OptimizedLIS finds the length of the longest increasing subsequence
// using an optimized approach with O(n log n) time complexity.
func OptimizedLIS(nums []int) int {
	// Check for empty slice in parameters
	length := len(nums)
	if length == 0 {
		return 0
	}

	// Initialization of helper slice op
	op := make([]int, 0)

	// For cycle to calculate helper slice with longest sequens
	for _, num := range nums {
		i := sort.Search(len(op), func(x int) bool {
			return op[x] >= num
		})
		if i == len(op) {
			op = append(op, num)
		} else {
			op[i] = num
		}
	}

	return len(op)
}

// GetLISElements returns one possible longest increasing subsequence
// (not just the length, but the actual elements).
func GetLISElements(nums []int) []int {
	// Check for empty slice in parameters
	length := len(nums)
	if length == 0 {
		return []int{}
	}
	// Initialization of helper slice op
	op := make([]int, 0)

	// For cyclle to calculate helper slice with longest sequense
	for _, num := range nums {
		i := sort.Search(len(op), func(x int) bool {
			return op[x] >= num
		})
		if i == len(op) {
			op = append(op, num)
		} else {
			op[i] = num
		}
	}

	return op
}
