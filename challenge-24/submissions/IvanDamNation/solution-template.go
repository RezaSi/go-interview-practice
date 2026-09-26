package main

import (
	"fmt"
	"slices"
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
	if len(nums) == 0 {
	    return 0
	}
	
	dp := make([]int, len(nums))
	for i := range dp {
	    dp[i] = 1
	}
	
	for i := range nums {
	    for j := 0; j < i; j++ {
	        if nums[i] > nums[j] && dp[i] < dp[j]+1 {
	            dp[i] = dp[j] + 1
	        }
	    }
	}
	
	maxVal := 1
	for _, v := range dp {
	    if maxVal < v {
	        maxVal = v
	    }
	}
	
	return maxVal
}

// OptimizedLIS finds the length of the longest increasing subsequence
// using an optimized approach with O(n log n) time complexity.
func OptimizedLIS(nums []int) int {
	if len(nums) == 0 {
	    return 0
	}
	
	tails := make([]int, 0, len(nums))
	for _, v := range nums {
	    idx, _ := slices.BinarySearch(tails, v)
	    if idx == len(tails) {
	        tails = append(tails, v)
	    } else {
	        tails[idx] = v
	    }
	}
	
	return len(tails)
}

// GetLISElements returns one possible longest increasing subsequence
// (not just the length, but the actual elements).
func GetLISElements(nums []int) []int {
	if len(nums) == 0 {
	    return []int{}
	}
	
	dp := make([]int, len(nums))
	for i := range dp {
	    dp[i] = 1
	}
	
	refs := make([]int, len(nums))
	for i := range refs {
	    refs[i] = -1
	}
	
	for i := range nums {
	    for j := 0; j < i; j++ {
	        if nums[i] > nums[j] && dp[i] < dp[j]+1 {
	            dp[i] = dp[j] + 1
	            refs[i] = j
	        }
	    } 
	}
	
	longest, longestIdx := 1, 0
	for i, v := range dp {
	    if longest < v {
	        longest = v
	        longestIdx = i
	    }
	}
	
	res := make([]int, longest)
	
	fillIdx := longest - 1
	currIdx := longestIdx
	for currIdx != -1 {
	    res[fillIdx] = nums[currIdx]
	    fillIdx--
	    currIdx = refs[currIdx]
	}
	
	return res
}
