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
	if amount == 0 {
		return 0
	}
	sort.Sort(sort.Reverse(sort.IntSlice(denominations)))
	remainigAmount := amount
	coinCount := 0

	for _, coin := range denominations {
		count := remainigAmount / coin
		if count == 0 {
			continue
		}
		coinCount += count
		remainigAmount -= coin * count
		if remainigAmount == 0 {
			return coinCount
		}
	}

	if remainigAmount != 0 {
		return -1
	}
	return 0
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	coinComb := make(map[int]int)
	if amount == 0 {
		return coinComb
	}
	sort.Sort(sort.Reverse(sort.IntSlice(denominations)))
	remainigAmount := amount

	for _, coin := range denominations {
		count := remainigAmount / coin
		if count == 0 {
			continue
		}
		coinComb[coin] = count
		remainigAmount -= coin * count
		if remainigAmount == 0 {
			return coinComb
		}
	}

	if remainigAmount != 0 {
		return coinComb
	}
	return nil
}
