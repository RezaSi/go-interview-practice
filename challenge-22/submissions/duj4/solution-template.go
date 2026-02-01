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
	// TODO: Implement this function
	// initialize a slice assigning each element a maximum value, except the 1st one
	dp := make([]int, amount + 1)
	for i := 1; i <= amount; i++ {
	    dp[i] = math.MaxInt
	}
	
	// calculate number of coins for each amount
	for _, coin := range denominations {
	    for i := coin; i <= amount; i++ {
	        if dp[i-coin] != math.MaxInt && dp[i-coin]+1 < dp[i] {
	            dp[i] = dp[i-coin] + 1
	        }
	    }
	}
	
	if dp[amount] == math.MaxInt {
	    return -1
	}
	
	return dp[amount]
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	// TODO: Implement this function
	dp := make([]int, amount+1)
	for i := 1; i <= amount; i++ {
	    dp[i] = math.MaxInt
	}
	

	
	// last coin used by the amount
	coinUsed := make([]int, amount+1)
	
	// mapping between amount and coint
	coinCombination := make(map[int]int)
	
	for _, coin := range denominations {
	    for i := coin; i <= amount; i++ {
	        if dp[i-coin] != math.MaxInt && dp[i-coin] + 1 < dp[i] {
	            dp[i] = dp[i-coin] + 1
	            coinUsed[i] = coin
	        }
	    }
	}
	
	if dp[amount] == math.MaxInt {
	    return coinCombination
	}
	
	for amount > 0 {
	    coin := coinUsed[amount]
	    coinCombination[coin]++
	    amount -= coin
	}

	return coinCombination
}
