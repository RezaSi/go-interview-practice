package main

import (
	"fmt"
	"slices"
)

func main() {
	// Standard U.S. coin denominations in cents
	denominations := []int{1, 5, 10, 25, 50}
	// denominations := []int{25, 10, 5}

	// Test amounts
	amounts := []int{87, 42, 99, 33, 7}
	// amounts := []int{3}

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
	slices.Sort(denominations)
	slices.Reverse(denominations)
	coinCount := 0
	for _, denomination := range denominations {
		if denomination <= amount {
			num := amount / denomination
			coinCount += num
			amount -= denomination * num
		}
	}

	if coinCount == 0 {
		return -1
	}

	return coinCount
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	combination := make(map[int]int, len(denominations))
	for _, denomination := range denominations {
		combination[denomination] = 0
	}
	slices.Sort(denominations)
	slices.Reverse(denominations)
	for _, denomination := range denominations {
		if denomination <= amount {
			num := amount / denomination
			combination[denomination] += num
			amount -= denomination * num
		}
	}
	for key, val := range combination {
		if val == 0 {
			delete(combination, key)
		}
	}

	return combination
}
