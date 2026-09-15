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

	dp := make([]int, len(nums), len(nums))

	for i := range len(dp) {
		dp[i] = 1
	}

	max := func(ints ...int) int {
		first := true
		var mx int
		for _, n := range ints {
			if first {
				mx = n
				first = false
			}
			if n > mx {
				mx = n
			}

		}
		return mx
	}

	for i := 0; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] {
				dp[i] = max(dp[i], dp[j]+1)
			}
		}
	}

	return max(dp...)
}

// OptimizedLIS finds the length of the longest increasing subsequence
// using an optimized approach with O(n log n) time complexity.
func OptimizedLIS(nums []int) int {

	tails := []int{}

	lowerBound := func(at []int, x int) int {
		l := 0
		r := len(at)

		for l < r {
			p := (l + r) / 2

			if at[p] < x {
				l = p + 1
			} else {
				r = p
			}

		}
		return l
	}

	for _, a := range nums {
		pos := lowerBound(tails, a)
		if pos == len(tails) {
			tails = append(tails, a)
		} else {
			tails[pos] = a
		}
	}

	return len(tails)
}

// GetLISElements returns one possible longest increasing subsequence
// (not just the length, but the actual elements).
func GetLISElements(nums []int) []int {
	if len(nums) == 0 {
		return nums
	}

	dp := make([]int, len(nums), len(nums))

	parent := make([]int, len(nums), len(nums))
	for i := range len(dp) {
		dp[i] = 1
		parent[i] = -1
	}

	/*
		max := func(ints ...int) int {
			first := true
			var mx int
			for _, n := range ints {
				if first {
					mx = n
					first = false
				}
				if n > mx {
					mx = n
				}

			}
			return mx
		}
	*/

	for i := 0; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
				parent[i] = j
			}
		}
	}

	//fmt.Println("parent", parent)
	//fmt.Println("dp", dp)

	idx := func(ints []int) int {
		first := true
		var mx, mi int
		for i, n := range ints {
			if first {
				mi = i
				first = false
			}
			if n > mx {
				mi = i
			}

		}
		return mi
	}(dp)

	elements := []int{}

	for {
		if idx == -1 {
			break
		}
		elements = append([]int{0}, elements...)
		elements[0] = nums[idx]
		idx = parent[idx]

	}

	return elements
}
