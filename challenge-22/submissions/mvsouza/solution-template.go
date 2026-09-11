package main

import (
	"fmt"
)

func main() {
	// Standard U.S. coin denominations in cents
	denominations := []int{1, 5, 10, 25, 50}

	// Test amounts
	amounts := []int{87, 42, 99, 33, 7}

	for _, amount := range amounts {
		// Find minimum number of coins
		minCoins := MinCoins(amount, denominations)

		// Find coin combination
		coinCombo := CoinCombination(amount, denominations)

		// Print results
		fmt.Printf("Amount: %d cents\n", amount)
		fmt.Printf("Minimum coins needed: %d\n", minCoins)
		fmt.Printf("Coin combination: %v\n", coinCombo)
		fmt.Println("---------------------------")
	}
}

// MinCoins returns the minimum number of coins needed to make the given amount.
// If the amount cannot be made with the given denominations, return -1.
func MinCoins(amount int, denominations []int) int {
	if amount == 0 {
		return 0
	}
	coinMap := CoinCombination(amount, denominations)
	if len(coinMap) == 0 {
		return -1
	}
	count := 0
	for _, quantity := range coinMap {
		count += quantity
	}
	return count
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	if amount < 0 {
		return map[int]int{}
	}
	const inf = int(^uint(0) >> 1)
	dp := make([]int, amount+1)
	prev := make([]int, amount+1) // chosen coin for reconstruction
	for i := 1; i <= amount; i++ {
		dp[i] = inf
		prev[i] = -1
	}
	for a := 1; a <= amount; a++ {
		for _, c := range denominations {
			if c <= 0 || c > a || dp[a-c] == inf {
				continue
			}
			if dp[a-c]+1 < dp[a] {
				dp[a] = dp[a-c] + 1
				prev[a] = c
			}
		}
	}
	if amount > 0 && prev[amount] == -1 {
		return map[int]int{}
	}
	coinMap := make(map[int]int)
	for a := amount; a > 0; {
		c := prev[a]
		if c <= 0 {
			return map[int]int{}
		}
		coinMap[c]++
		a -= c
	}
	return coinMap
}
