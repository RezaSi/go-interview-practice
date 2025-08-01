package main

import (
	"fmt"
)

const inf = 2000000000

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
	var LIS = make([]int, len(nums))
	answer := 0

	for i, num := range nums {
		LIS[i] = 1
		for j := 0; j < i; j++ {
			if nums[j] < num {
				LIS[i] = max(LIS[i], LIS[j]+1)
			}
		}
		answer = max(answer, LIS[i])
	}

	return answer
}

func find(a *[]int, val, ind int) int {
	L, R := 0, ind
	var M int

	for R-L > 1 {
		M = (L + R) / 2
		if (*a)[M] < val {
			L = M
		} else {
			R = M
		}
	}

	return L
}

// OptimizedLIS finds the length of the longest increasing subsequence
// using an optimized approach with O(n log n) time complexity.
func OptimizedLIS(nums []int) int {
	answer := 0

	minVal := make([]int, len(nums)+1)
	minVal[0] = -inf
	for i := 1; i <= len(nums); i++ {
		minVal[i] = inf
	}

	var ind int

	for i, num := range nums {
		ind = find(&minVal, num, i+1)
		answer = max(answer, ind+1)
		minVal[ind+1] = min(minVal[ind+1], num)
	}

	return answer
}

func buildLIS(LIS, a, p *[]int, index int) {
	if index == -1 {
		return
	}

	buildLIS(LIS, a, p, (*p)[index])
	*LIS = append(*LIS, (*a)[index])
}

// GetLISElements returns one possible longest increasing subsequence
// (not just the length, but the actual elements).
func GetLISElements(nums []int) []int {
	answer, bestInd := 0, -1
	minVal := make([]int, len(nums)+1)
	minInd := make([]int, len(nums)+1)
	pre := make([]int, len(nums))

	minVal[0] = -inf
	minInd[0] = -1
	for i := 1; i <= len(nums); i++ {
		minVal[i] = inf
	}

	var ind int

	for i, num := range nums {
		ind = find(&minVal, num, i+1)
		if ind+1 > answer {
			answer = ind + 1
			bestInd = i
		}
		if num < minVal[ind+1] {
			minVal[ind+1] = num
			minInd[ind+1] = i
		}
		pre[i] = minInd[ind]
	}

	seq := make([]int, 0)
	buildLIS(&seq, &nums, &pre, bestInd)
	return seq
}
