package main

import (
	"fmt"
	"sort"
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
	sort.Sort(sort.Reverse(sort.IntSlice(denominations)))
	count := 0
	for _, coin := range denominations {
		coinCount := amount / coin
		count += coinCount
		amount -= coin * coinCount
		if amount == 0 {
			return count
		}
	}

	if amount > 0 {
		return -1
	}
	return count
}

func CoinCombination(amount int, denominations []int) map[int]int {
	sort.Sort(sort.Reverse(sort.IntSlice(denominations)))
	result := make(map[int]int)

	for _, coin := range denominations {
		coinCount := amount / coin
		amount -= coin * coinCount
		if coinCount > 0 {
		    result[coin] = coinCount
		}
		if amount == 0 {
			return result
		}
	}

	return result
}
