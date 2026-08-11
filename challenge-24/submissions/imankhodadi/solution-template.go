package main

import (
	"fmt"
	"sort"
)

func main() {
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

func DPLongestIncreasingSubsequence(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	//dp[i] represents the length of the longest increasing subsequence ending at index i.
	dp := make([]int, len(nums))
	for i := range dp {
		dp[i] = 1 // Every element is a subsequence of length 1
	}
	for i := 1; i < len(nums); i++ {
		//check all previous elements that are smaller
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] {
				if dp[j]+1 > dp[i] {
					dp[i] = dp[j] + 1
				}
			}
		}
	}
	max := dp[0]
	for _, x := range dp {
		if x > max {
			max = x
		}
	}
	return max
}

func OptimizedLIS(nums []int) int {
	tails := []int{} //tails[i] represents the smallest value at which an increasing subsequence of length i+1 ends
	for _, num := range nums {
		pos := sort.SearchInts(tails, num)
		if pos == len(tails) {
			tails = append(tails, num) // Extend the sequence, If larger than all elements in tails, append it
		} else {
			tails[pos] = num // Replace with smaller ending element, Otherwise, find the first element >= current element and replace it
		}
	}
	return len(tails) //The length of tails array is the LIS length
}

func GetLISElements(nums []int) []int {
	if len(nums) == 0 {
		return []int{}
	}
	dp := make([]int, len(nums))
	for i := range dp {
		dp[i] = 1 
	}
	parent := make([]int, len(nums))
	for i := 1; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] {
				if dp[j]+1 > dp[i] {
					dp[i] = dp[j] + 1
					parent[i] = j
				}
			}
		}
	}
	var maxIndex int
	maxVal := dp[0]
	for i, x := range dp {
		if x > maxVal {
			maxVal, maxIndex = x, i
		}
	}
	lis := make([]int, maxVal)
	current := maxIndex
	for i := maxVal - 1; i >= 0; i-- {
		lis[i] = nums[current]
		current = parent[current]
	}
	return lis
}