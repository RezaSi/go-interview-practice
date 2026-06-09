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
	if amount < 0 || len(denominations) == 0 {
		return -1
	}

	for _, d := range denominations {
		if d <= 0 {
			return -1
		}
	}

	const inf = int(^uint(0) >> 1)
	dp := make([]int, amount+1)

	for i := 1; i <= amount; i++ {
		dp[i] = inf
	}

	for a := 1; a <= amount; a++ {
		for _, d := range denominations {
			if d <= a && dp[a-d] != inf && dp[a-d]+1 < dp[a] {
				dp[a] = dp[a-d] + 1
			}
		}
	}
	if dp[amount] == inf {
		return -1
	}
	return dp[amount]
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	if amount <= 0 || len(denominations) == 0 {
		return map[int]int{}
	}
	for _, d := range denominations {
		if d <= 0 {
			return map[int]int{}
		}
	}

	const inf = int(^uint(0) >> 1)
	dp := make([]int, amount+1)
	pick := make([]int, amount+1)
	for i := 1; i <= amount; i++ {
		dp[i] = inf
		pick[i] = -1
	}

	for a := 1; a <= amount; a++ {
		for _, d := range denominations {
			if d <= a && dp[a-d] != inf && dp[a-d]+1 < dp[a] {
				dp[a] = dp[a-d] + 1
				pick[a] = d
			}
		}
	}
	if dp[amount] == inf {
		return map[int]int{}
	}

	res := map[int]int{}
	for a := amount; a > 0; {
		d := pick[a]
		if d <= 0 {
			return map[int]int{}
		}
		res[d]++
		a -= d
	}
	return res
}
