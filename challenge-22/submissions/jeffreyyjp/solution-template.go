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
	// TODO: Implement this function
	max := amount + 1
	dp := make([]int, max)
	for i := 1; i < len(dp); i++ {
	    dp[i] = max
	}
	
	for i := 1; i <= amount; i++ {
	    for _, coin := range denominations {
	        if i >= coin {
	            dp[i] = min(dp[i], dp[i-coin] + 1)
	        }
	    }
	}
	
	if dp[amount] >= max {
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
	if amount == 0 || len(denominations) == 0 {
	    return make(map[int]int)
	}

	max := amount + 1
	dp := make([]int, max)
	combinations := make([]map[int]int, max)
	for i := 1; i < len(dp); i++ {
	    dp[i] = max
	    combinations[i] = make(map[int]int)
	}
	combinations[0] = make(map[int]int)
	
	for i := 1; i <= amount; i++ {
	    for _, coin := range denominations {
	        if i >= coin && dp[i] > dp[i-coin] + 1 {
	            dp[i] = dp[i-coin] + 1
	            
	            combinations[i] = make(map[int]int)
	            for k, v := range combinations[i-coin] {
	                combinations[i][k] = v
	            }
	            combinations[i][coin]++
	        }
	    }
	}
	
	if dp[amount] >= max {
	    return make(map[int]int)
	}
	
	return combinations[amount]
}
