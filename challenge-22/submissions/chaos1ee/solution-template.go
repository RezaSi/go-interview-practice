package main

import "sort"

// MinCoins returns the minimum number of coins needed to make the given amount.
// Uses greedy approach: always choose the largest coin possible.
// If the amount cannot be made with the given denominations, return -1.
func MinCoins(amount int, denominations []int) int {
	if amount == 0 {
		return 0
	}

	denom := make([]int, len(denominations))
	copy(denom, denominations)
	sort.Sort(sort.Reverse(sort.IntSlice(denom)))

	coinCount := 0
	remainingAmount := amount

	for _, coin := range denom {
		count := remainingAmount / coin
		coinCount += count
		remainingAmount -= count * coin

		if remainingAmount == 0 {
			return coinCount
		}
	}

	return -1
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// Uses greedy approach. If the amount cannot be made, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	if amount == 0 {
		return map[int]int{}
	}

	denom := make([]int, len(denominations))
	copy(denom, denominations)
	sort.Sort(sort.Reverse(sort.IntSlice(denom)))

	result := make(map[int]int)
	remainingAmount := amount

	for _, coin := range denom {
		count := remainingAmount / coin
		remainingAmount -= coin * count
		if count > 0 {
			result[coin] = count
		}

		if remainingAmount == 0 {
			return result
		}
	}

	if remainingAmount > 0 {
		return map[int]int{}
	}
	return result
}
