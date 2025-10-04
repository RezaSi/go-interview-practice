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
	n := len(denominations)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, amount+1)
	}
	for i := 0; i <= amount; i++ {
		dp[0][i] = math.MaxInt - 1
	}
	for i := 0; i <= n; i++ {
		dp[i][0] = 0
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= amount; j++ {
			if denominations[i-1] <= j {
				dp[i][j] = min(dp[i][j-denominations[i-1]]+1, dp[i-1][j])
			} else {
				dp[i][j] = dp[i-1][j]
			}
		}
	}
    if dp[n][amount] >= math.MaxInt - 1 {
        return -1
    }
    
	return dp[n][amount]
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	// TODO: Implement this function
	n := len(denominations)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, amount+1)
	}
	for i := 0; i <= amount; i++ {
		dp[0][i] = math.MaxInt - 1
	}
	for i := 0; i <= n; i++ {
		dp[i][0] = 0
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= amount; j++ {
			if denominations[i-1] <= j {
				dp[i][j] = min(dp[i][j-denominations[i-1]]+1, dp[i-1][j])
			} else {
				dp[i][j] = dp[i-1][j]
			}
		}
	}
	ans := map[int]int{}

    if dp[n][amount] >= math.MaxInt - 1 {
        return ans
    }
    i := n
	j := amount

    for i > 0 && j > 0 {
        coin := denominations[i-1]
        
        if j >= coin && dp[i][j] == dp[i][j-coin]+1 {
            ans[coin]++
            j -= coin
        } else {
            i--
        }
    }
    
	return ans
}
