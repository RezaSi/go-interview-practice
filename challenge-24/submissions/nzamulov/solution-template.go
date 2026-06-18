package main

import (
	"fmt"
	"math"
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
    result := 0
	dp := make([]int, len(nums))
	for i := range nums {
	    dp[i] = 1
	    for j := 0; j < i; j++ {
	        if nums[j] < nums[i] {
	            dp[i] = max(dp[i], dp[j] + 1)
	        }
	    }
	    result = max(result, dp[i])
	}
	return result
}

// OptimizedLIS finds the length of the longest increasing subsequence
// using an optimized approach with O(n log n) time complexity.
func OptimizedLIS(nums []int) int {
	return len(GetLISElements(nums))
}

// GetLISElements returns one possible longest increasing subsequence
// (not just the length, but the actual elements).
func GetLISElements(nums []int) []int {
    n := len(nums)
	dp := make([]int, n+1)
	pos := make([]int, n+1)
	prev := make([]int, n+1)
	for i := 0; i < n+1; i++ {
	    if i == 0 {
	        dp[i] = math.MinInt32
	    } else {
	        dp[i] = math.MaxInt32
	    }
	    pos[i] = -1
	}
	length := 0
	for i := 0; i < n; i++ {
	    l, r := 0, n
	    for l + 1 < r {
	        m := (l + r) >> 1
	        if dp[m] >= nums[i] {
	            r = m
	        } else {
	            l = m
	        }
	    }
        dp[r] = nums[i]
	    pos[r] = i
	    prev[i] = pos[r - 1]
	    length = max(length, r)
	}
	var result []int
	if n > 0 {
	    p := pos[length]
    	for p != -1 {
    	    result = append(result, nums[p])
    	    p = prev[p]
    	}
	}
	for i := 0; i < length / 2; i++ {
	    result[i], result[length - i - 1] = result[length - i - 1], result[i]
	}
	return result
}
