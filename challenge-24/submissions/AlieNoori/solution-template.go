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
	if n == 0 {
		return 0
	}

	dp := make([]int, n)
	for i := range dp {
		dp[i] = 1
	}

	longest := 1
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
			}
		}
		if dp[i] > longest {
			longest = dp[i]
		}
	}

	return longest
}

// OptimizedLIS finds the length of the longest increasing subsequence
// using an optimized approach with O(n log n) time complexity.
func OptimizedLIS(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}

	tails := []int{}

	for i := 0; i < n; i++ {
		idx := findInsertIndex(tails, nums[i])

		if idx == len(tails) {
			tails = append(tails, nums[i])
		} else {
			tails[idx] = nums[i]
		}
	}

	return len(tails)
}

func findInsertIndex(nums []int, target int) int {
	if len(nums) == 0 {
		return 0
	}

	low, high := 0, len(nums)

	for high > low {
		mid := (high-low)/2 + low
		if target > nums[mid] {
			low = mid + 1
		} else {
			high = mid
		}
	}

	return low
}

// GetLISElements returns one possible longest increasing subsequence
// (not just the length, but the actual elements).
func OptimizedLISElements(nums []int) []int {
	n := len(nums)
	if n == 0 {
		return []int{}
	}

	prev := make([]int, n)
	tailsIdx := []int{}
	tails := []int{}

	for i, x := range nums {
		idx := findInsertIndex(tails, x)

		if idx > 0 {
			prev[i] = tailsIdx[idx-1]
		} else {
			prev[i] = -1
		}

		if idx == len(tails) {
			tails = append(tails, x)
			tailsIdx = append(tailsIdx, i)
		} else {
			tails[idx] = x
			tailsIdx[idx] = i
		}
	}

	length := len(tails)
	result := make([]int, length)
	k := tailsIdx[length-1]
	for i := length - 1; i >= 0; i-- {
		result[i] = nums[k]
		k = prev[k]
	}

	return result
}

// GetLISElements returns one possible longest increasing subsequence
// (not just the length, but the actual elements).
func GetLISElements(nums []int) []int {
	n := len(nums)
	if n == 0 {
		return []int{}
	}

	dp := make([]int, n)
	prev := make([]int, n)
	for i := range dp {
		dp[i] = 1
		prev[i] = -1
	}

	best := 0
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
				prev[i] = j
			}
		}

		if dp[i] > dp[best] {
			best = i
		}
	}

	length := dp[best]
	result := make([]int, length)
	for k := best; k != -1; k = prev[k] {
		length--
		result[length] = nums[k]
	}

	return result
}
