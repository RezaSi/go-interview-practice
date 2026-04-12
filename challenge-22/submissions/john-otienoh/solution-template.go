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
    dp := make([]int, amount+1)
    choice := make([]int, amount+1)
    for i := 1; i <= amount; i++ {
        dp[i] = 1 << 30
    }
    dp[0] = 0

    for i := 1; i <= amount; i++ {
        for _, coin := range denominations {
            if i-coin >= 0 && dp[i-coin]+1 < dp[i] {
                dp[i] = dp[i-coin] + 1
                choice[i] = coin
            }
        }
    }
    if dp[amount] == 1<<30 {
        return -1
    }

    result := []int{}
    curr := amount
    for curr > 0 {
        c := choice[curr]
        result = append(result, c)
        curr -= c
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
    choice := make([]int, amount+1)
    for i := 1; i <= amount; i++ {
        dp[i] = 1 << 30
    }
    dp[0] = 0

    for i := 1; i <= amount; i++ {
        for _, coin := range denominations {
            if i-coin >= 0 && dp[i-coin]+1 < dp[i] {
                dp[i] = dp[i-coin] + 1
                choice[i] = coin
            }
        }
    }
    if dp[amount] == 1<<30 {
        return map[int]int{}
    }
    combination := make(map[int]int)
    curr := amount
    for curr > 0 {
        coin := choice[curr]
        combination[coin]++
        curr -= coin
    }

    return combination
	
}
