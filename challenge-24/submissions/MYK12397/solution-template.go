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
	// TODO: Implement this function
    if len(nums) == 0 {
        return 0
    }
    n := len(nums)
    if n == 1 {
        return 1
    }
    
    dp := make([]int,n+1)
    dp[0] = 1
    res := -1

    for i := 1; i < n; i++ {
        mx := 0
        for j := 0; j < i; j++ {
            if nums[i] > nums[j] {
                if mx <= dp[j] {
                    mx =  dp[j]
                }
        }
        dp[i]= mx + 1
        if res < dp[i] {
            res = dp[i]
        }
    }
    }
    return res
}

// OptimizedLIS finds the length of the longest increasing subsequence
// using an optimized approach with O(n log n) time complexity.
func OptimizedLIS(nums []int) int {
	// TODO: Implement this function
	var tails []int
	for i := 0; i < len(nums); i++ {
		tlen := len(tails)
		if tlen == 0 {
			tails = []int{nums[i]}
			continue
		}
		if nums[i] > tails[tlen-1] {
			tails = append(tails, nums[i])
			continue
		}

		start, end, mid := 0, tlen-1, 0
		for start != end {
			mid = start + (end-start)/2
			if nums[i] > tails[mid] {
				start = mid + 1
			} else {
				end = mid
			}
		}
		tails[end] = nums[i]
	}
	return len(tails)
}

// GetLISElements returns one possible longest increasing subsequence
// (not just the length, but the actual elements).
func GetLISElements(nums []int) []int {
	// TODO: Implement this function
	n := len(nums)
	if n == 0 {
		return []int{}
	}
	
	dp := make([]int, n)
	parent := make([]int, n)
	
	for i := 0; i < n; i++ {
		dp[i] = 1
		parent[i] = -1
	}
	
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
				parent[i] = j
			}
		}
	}
	maxLength := dp[0]
	maxIndex := 0
	for i := 1; i < n; i++ {
		if dp[i] > maxLength {
			maxLength = dp[i]
			maxIndex = i
		}
	}
	
	result := make([]int, maxLength)
	current := maxIndex
	for i := maxLength - 1; i >= 0; i-- {
		result[i] = nums[current]
		current = parent[current]
	}
	
	return result
}
