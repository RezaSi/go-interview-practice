package main

import (
	"fmt"
	"slices"
)

func MinCoins(amount int, denominations []int) int {
	if amount < 0 {
		// fmt.Errorf("Cannot take negative amount") use -1
		return -1//,nil
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
		// fmt.Errorf("Cannot give this amount with these coins")
		return -1//,nil
	}
	return coinsCount
}

func CoinCombination(amount int, denominations []int) map[int]int {
	if amount < 0 {
		// fmt.Errorf("Cannot take negative amount")
		return map[int]int{}//,nil
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
		// fmt.Errorf("Cannot give this amount with these coins")
		return map[int]int{}//,nil
	}
	return coinsMap
}

func MinCoinsDP(amount int, denominations []int) int {
	if amount < 0 {
		// fmt.Errorf("Cannot take negative amount")
		return 0//,nil
	}
	if amount == 0 {
		return 0
	}
	dp := make([]int, amount+1)
	for i := range dp {
		dp[i] = amount + 1
	}
	dp[0] = 0 // Base case: 0 coins needed to make amount 0
	for _, coin := range denominations {
		for i := coin; i <= amount; i++ {
			if dp[i-coin]+1 < dp[i] {
				dp[i] = dp[i-coin] + 1
			}
		}
	}
	if dp[amount] > amount {
		return -1//,nil
	}
	return dp[amount]
}

func main() {
	denominations := []int{1, 5, 10, 25, 50}
	for _, x := range []int{0, -1, 20} {
		fmt.Println(MinCoins(x, denominations),
			CoinCombination(x, denominations),
			MinCoinsDP(x, denominations))
	}

	denominations = []int{3, 5}
	for _, x := range []int{0, -1, 20, 2} {
		fmt.Println(MinCoins(x, denominations),
			CoinCombination(x, denominations),
			MinCoinsDP(x, denominations))
	}
	denominations = []int{1, 3, 4}
	for _, x := range []int{0, -1, 20, 2, 6} {
		fmt.Println(MinCoins(x, denominations),
			CoinCombination(x, denominations),
			MinCoinsDP(x, denominations))
	}
}