package main

import (
	"fmt"
	"sort"
)

func main() {
	// Standard U.S. coin denominations in cents
	denominations := []int{1, 5, 10, 25, 50}

	// Test amounts
	amounts := []int{0}

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
func MinCoins(amount int, denominations []int) (c int) {
	if amount == 0 {
		return amount
	}
	coins := CoinCombination(amount, denominations)
	if len(coins) == 0 {
		return -1
	}
	for _, coin := range coins {
		c += coin
	}

	return
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	result := map[int]int{}
	var count int
	sort.Slice(denominations, func(i, j int) bool {
		return denominations[i] > denominations[j]
	})
	for _, v := range denominations {
		for amount-v >= 0 {
			amount -= v
			result[v] += 1
			count++
		}
	}

	return result
}
