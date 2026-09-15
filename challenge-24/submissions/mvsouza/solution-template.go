package main

import (
	"fmt"
	"sync"
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

var increasingSubsequenceMap sync.Map

// DPLongestIncreasingSubsequence finds the length of the longest increasing subsequence
// using a standard dynamic programming approach with O(n²) time complexity.
func DPLongestIncreasingSubsequence(nums []int) int {
	highestValue := 1
	arr := make([]int, len(nums))
	if len(nums) == 0 {
		return 0
	}
	arr[0] = 1
	for i := 1; i < len(nums); i++ {
		arr[i] = 1
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] && arr[j] >= arr[i] {
				arr[i] = arr[j] + 1
			}
		}
		if arr[i] > highestValue {
			highestValue = arr[i]
		}
	}
	return highestValue
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
	if n == 0 {
		return []int{}
	}
	tailsIdx := make([]int, 0, n)
	prev := make([]int, n)
	for i := range prev {
		prev[i] = -1
	}
	for i, num := range nums {
		l, r := 0, len(tailsIdx)
		for l < r {
			m := l + (r-l)/2
			if nums[tailsIdx[m]] < num {
				l = m + 1
			} else {
				r = m
			}
		}
		if l > 0 {
			prev[i] = tailsIdx[l-1]
		}
		if l == len(tailsIdx) {
			tailsIdx = append(tailsIdx, i)
		} else {
			tailsIdx[l] = i
		}
	}
	k := tailsIdx[len(tailsIdx)-1]
	res := make([]int, len(tailsIdx))
	for p := len(res) - 1; p >= 0; p-- {
		res[p] = nums[k]
		k = prev[k]
	}
	return res
}
