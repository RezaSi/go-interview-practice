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
	if amount < 1 {
	    return 0
	}
	dp := make([]int,amount+1)
	for i := range dp {
	    dp[i] = math.MaxInt32
	}
	dp[0] = 0
	for i := 1;i< len(dp);i++ {
	    for _,val := range denominations {
	        if val <= i {
	            rem := i -val
	            dp[i] = min(dp[i],dp[rem] + 1)
	        }
	    }
	}
	if dp[amount] == math.MaxInt32 {
	    return -1
	}
	return dp[amount]
}

func min(i,j int) int {
    if i < j {
        return i
    }
    return j
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	// TODO: Implement this function
	if amount < 1 {
	    return map[int]int{}
	}
	dp := make([]int,amount+1)
	prev := make(map[int]int)
	for i := range dp {
	    dp[i] = math.MaxInt32
	}
	dp[0] = 0
	for i := 1;i< len(dp);i++ {
	    for _,val := range denominations {
	        if val <= i {
	            rem := i - val
	            if dp[rem] + 1 < dp[i] {
	                dp[i] = dp[rem] + 1
	                prev[i] = val
	            }
	            
	        }
	    }
	}
	coinComb := make(map[int]int)
	if dp[amount] == math.MaxInt32 {
		return coinComb
	}
	for amt := amount; amt > 0; amt -= prev[amt] {
		coinComb[prev[amt]]++
	}
	return coinComb
}
