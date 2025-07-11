package main

import (
	"fmt"
	"math"
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
	const MaxInt = math.MaxInt32

	dp := make([]int, amount+1)
	for i := range dp {
		dp[i] = MaxInt
	}
	dp[0] = 0

	for am := 1; am <= amount; am++ {
		for _, coin := range denominations {
			if am >= coin && dp[am-coin] != MaxInt {
				dp[am] = min(dp[am], dp[am-coin]+1)
			}
		}
	}

	if dp[amount] == MaxInt {
		return -1
	}
	return dp[amount]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}


// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	const MaxInt = int(^uint(0) >> 1)

	dp := make([]int, amount+1)
	lastCoin := make([]int, amount+1)

	for i := range dp {
		dp[i] = MaxInt
	}
	dp[0] = 0

	for am := 1; am <= amount; am++ {
		for _, coin := range denominations {
			if am >= coin && dp[am-coin] != MaxInt {
				if dp[am-coin]+1 < dp[am] {
					dp[am] = dp[am-coin] + 1
					lastCoin[am] = coin
				}
			}
		}
	}

	if dp[amount] == MaxInt {
		return map[int]int{}
	}

	result := make(map[int]int)
	for amt := amount; amt > 0; {
		coin := lastCoin[amt]
		result[coin]++
		amt -= coin
	}

	return result
}

