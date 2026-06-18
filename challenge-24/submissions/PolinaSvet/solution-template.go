package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	// Test cases
	testCases := []struct {
		name string
		nums []int
	}{
		{"Example 01: 4", []int{10, 9, 2, 5, 3, 7, 101, 18}},
		{"Example 02: 4", []int{0, 1, 0, 3, 2, 3}},
		{"Example 03: 1", []int{7, 7, 7, 7, 7, 7, 7}},
		{"Example 04: 3", []int{4, 10, 4, 3, 8, 9}},
		{"Example 05: 0", []int{}},
		{"Example 06: 1", []int{5}},
		{"Example 07: 1", []int{5, 4, 3, 2, 1}},
		{"Example 08: 5", []int{1, 2, 3, 4, 5}},
		{"Example 09: 3", []int{3, 10, 2, 1, 20}},
		{"Example 10: 4", []int{50, 3, 10, 7, 40, 80}},
	}

	// Test each approach
	for _, tc := range testCases {
		dpLength := DPLongestIncreasingSubsequence(tc.nums)
		fmt.Printf("Standart%s = %d, %v\n", tc.name, dpLength, tc.nums)

		optLength := OptimizedLIS(tc.nums)
		fmt.Printf("OPT_____%s = %d, %v\n", tc.name, optLength, tc.nums)

		lisElements := GetLISElements(tc.nums)
		fmt.Printf("LIS Elements: %v\n", lisElements)
		fmt.Println("-----------------------------------")
	}

	fmt.Println("\n=== Performance Tests ===")
	runPerformanceTests()

	fmt.Println("\n=== Memory Usage Tests ===")
	runMemoryTests()
}

// DPLongestIncreasingSubsequence finds the length of the longest increasing subsequence
// using a standard dynamic programming approach with O(n²) time complexity.
func DPLongestIncreasingSubsequence(nums []int) int {
	// TODO: Implement this function
	if len(nums) == 0 {
		return 0
	}

	dp := make([]int, len(nums))
	for i := range dp {
		dp[i] = 1
	}

	maxLength := 1
	for i := 1; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] {
				if dp[j]+1 > dp[i] {
					dp[i] = dp[j] + 1
				}
			}
		}
		if dp[i] > maxLength {
			maxLength = dp[i]
		}
	}

	return maxLength
}

// OptimizedLIS finds the length of the longest increasing subsequence
// using an optimized approach with O(n log n) time complexity.
func OptimizedLIS(nums []int) int {
	// TODO: Implement this function
	if len(nums) == 0 {
		return 0
	}

	tails := make([]int, 0)
	tails = append(tails, nums[0])

	for i := 1; i < len(nums); i++ {
		if nums[i] > tails[len(tails)-1] {
			tails = append(tails, nums[i])
		} else {

			left, right := 0, len(tails)-1
			for left < right {
				mid := left + (right-left)/2
				if tails[mid] < nums[i] {
					left = mid + 1
				} else {
					right = mid
				}
			}
			tails[left] = nums[i]
		}
	}

	return len(tails)
}

// GetLISElements returns one possible longest increasing subsequence
// (not just the length, but the actual elements).
func GetLISElements(nums []int) []int {
	// TODO: Implement this function
	if len(nums) == 0 {
		return nil
	}

	dp := make([]int, len(nums))
	prev := make([]int, len(nums))
	for i := range dp {
		dp[i] = 1
		prev[i] = -1
	}

	maxLength := 1
	maxIndex := 0

	for i := 1; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
				prev[i] = j
			}
		}
		if dp[i] > maxLength {
			maxLength = dp[i]
			maxIndex = i
		}
	}

	lis := make([]int, maxLength)
	for i := maxLength - 1; i >= 0; i-- {
		lis[i] = nums[maxIndex]
		maxIndex = prev[maxIndex]
		if maxIndex == -1 {
			break
		}
	}

	return lis
}

func generateTestData(size int) []int {
	rand.Seed(time.Now().UnixNano())
	data := make([]int, size)
	for i := range data {
		data[i] = rand.Intn(size * 2)
	}
	return data
}

func runPerformanceTests() {
	sizes := []int{100, 500, 1000, 2000, 5000}
	iterations := 10

	for _, size := range sizes {
		fmt.Printf("\nTesting with size %d (%d iterations):\n", size, iterations)
		testData := generateTestData(size)

		var totalTime1, totalTime2 time.Duration

		for i := 0; i < iterations; i++ {
			start := time.Now()
			DPLongestIncreasingSubsequence(testData)
			totalTime1 += time.Since(start)
		}

		for i := 0; i < iterations; i++ {
			start := time.Now()
			OptimizedLIS(testData)
			totalTime2 += time.Since(start)
		}

		avgTime1 := totalTime1 / time.Duration(iterations)
		avgTime2 := totalTime2 / time.Duration(iterations)

		fmt.Printf("Standart: avg time: %v\n", avgTime1)
		fmt.Printf("OPT_____: avg time: %v\n", avgTime2)

	}
}

// Memory test function
func runMemoryTests() {
	sizes := []int{100, 500, 1000, 2000}

	for _, size := range sizes {
		fmt.Printf("\nMemory test with size %d:\n", size)
		testData := generateTestData(size)

		var m1, m2 runtime.MemStats

		runtime.GC()
		runtime.ReadMemStats(&m1)
		mResult1 := DPLongestIncreasingSubsequence(testData)
		runtime.ReadMemStats(&m2)
		mMemory1 := m2.Alloc - m1.Alloc

		runtime.GC()
		runtime.ReadMemStats(&m1)
		mResult2 := OptimizedLIS(testData)
		runtime.ReadMemStats(&m2)
		mMemory2 := m2.Alloc - m1.Alloc

		fmt.Printf("Standart: result=%d, memory=%d bytes\n", mResult1, mMemory1)
		fmt.Printf("OPT_____: result=%d, memory=%d bytes\n", mResult2, mMemory2)

	}
}
