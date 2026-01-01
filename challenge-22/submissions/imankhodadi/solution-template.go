package main

import (
	"fmt"
	"slices"
)

func MinCoins(amount int, denominations []int) int {
	if amount < 0 {
		return -1
	}
	if amount == 0 {
		return 0
	}
	var res int
	sorted := make([]int, len(denominations))
	copy(sorted, denominations)
	slices.Sort(sorted)

	coinsCount := 0
	for i := len(sorted) - 1; i >= 0; i-- {
		res = amount / sorted[i]
		if res > 0 {
			coinsCount += res
			amount %= sorted[i]
		}
		if amount == 0 {
			break
		}
	}
	if amount > 0 {
		return -1
	}
	return coinsCount
}

func CoinCombination(amount int, denominations []int) map[int]int {
	if amount < 0 {
		return map[int]int{}
	}
	if amount == 0 {
		return map[int]int{}
	}
	var res int
	sorted := make([]int, len(denominations))
	copy(sorted, denominations)
	slices.Sort(sorted)
	coinsMap := map[int]int{}
	for i := len(sorted) - 1; i >= 0; i-- {
		res = amount / sorted[i]
		if res > 0 {
			coinsMap[sorted[i]] += res
			amount %= sorted[i]
		}
		if amount == 0 {
			break
		}
	}
	if amount > 0 {
		return map[int]int{}
	}
	return coinsMap
}
func MinCoinsDP(amount int, denominations []int) int {
	if amount < 0 {
		return 0
	}
	if amount == 0 {
		return 0
	}
	dp := make([]int, amount+1)
	for i := range dp {
		dp[i] = amount + 1
	}
	dp[0] = 0
	for _, coin := range denominations {
		for i := coin; i <= amount; i++ {
			if dp[i-coin]+1 < dp[i] {
				dp[i] = dp[i-coin] + 1
			}
		}
	}
	if dp[amount] > amount {
		return -1
	}
	return dp[amount]
}

func main() {
	// Implementations of 3 functions for give change problem. 
	// Two greedy functions are similar, but their outputs are different.
	// DP approach is optimized for cases when greedy fails.
	// all three will fail in cases when coins are not appropriate for the task
	// valid coins
	denominations := []int{1, 5, 10, 25, 50}
	for _, x := range []int{0, -1, 20} {
		fmt.Println(MinCoins(x, denominations),
			CoinCombination(x, denominations),
			MinCoinsDP(x, denominations))
	}
	// this example shows we cannot give change for some values like 2 with coins {3, 5}
	denominations = []int{3, 5}
	for _, x := range []int{0, -1, 20, 2} {
		fmt.Println(MinCoins(x, denominations),
			CoinCombination(x, denominations),
			MinCoinsDP(x, denominations))
	}
	// this example shows the greedy approach fails in some cases
	denominations = []int{1, 3, 4}
	for _, x := range []int{0, -1, 20, 2, 6} {
		fmt.Println(MinCoins(x, denominations),
			CoinCombination(x, denominations),
			MinCoinsDP(x, denominations))
	}
}