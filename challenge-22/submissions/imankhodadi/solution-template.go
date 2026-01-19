package main

import (
	"fmt"
	"slices"
)

func MinCoins(amount int, denominations []int) int {
	if amount < 0 {
		fmt.Println("Cannot take negative amounts")
		return -1
	}
	if amount == 0 {
		return 0
	}
	var res int
	slices.Sort(denominations)
	coinsCount := 0
	for i := len(denominations) - 1; i >= 0; i-- {
		res = amount / denominations[i]
		if res > 0 {
			coinsCount += res
			amount %= denominations[i]
		}
		if amount == 0 {
			break
		}
	}
	if amount > 0 {
		fmt.Println("Cannot give this amounts with these coins")
		return -1
	}
	return coinsCount
}

func CoinCombination(amount int, denominations []int) map[int]int {
	if amount < 0 {
		fmt.Println("Cannot take negative amounts")
		return map[int]int{}
	}
	if amount == 0 {
		return map[int]int{}
	}
	var res int
	slices.Sort(denominations)
	coinsMap := map[int]int{}
	for i := len(denominations) - 1; i >= 0; i-- {
		res = amount / denominations[i]
		if res > 0 {
			coinsMap[denominations[i]] += res
			amount %= denominations[i]
		}
		if amount == 0 {
			break
		}
	}
	if amount > 0 {
		fmt.Println("Cannot give this amounts with these coins")
		return map[int]int{}
	}
	return coinsMap
}

func minCoinsDP(amount int, denominations []int) int {
	if amount < 0 {
		fmt.Println("Cannot take negative amounts")
		return -1
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
		return -1
	}
	return dp[amount]
}

func main() {
	denominations := []int{1, 5, 10, 25, 50}
	for _, x := range []int{0, -1, 20} {
		fmt.Println(MinCoins(x, denominations),
			CoinCombination(x, denominations),
			minCoinsDP(x, denominations))
	}

	denominations = []int{3, 5}
	for _, x := range []int{0, -1, 20, 2} {
		fmt.Println(MinCoins(x, denominations),
			CoinCombination(x, denominations),
			minCoinsDP(x, denominations))
	}
}