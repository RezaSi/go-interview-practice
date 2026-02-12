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

func reverse(slice []int) []int {
	newSlice := make([]int, len(slice), len(slice))
	copy(newSlice, slice)

	for i, j := 0, len(newSlice)-1; i < j; i, j = i+1, j-1 {
		newSlice[i], newSlice[j] = newSlice[j], newSlice[i]
	}

	return newSlice
}

// MinCoins returns the minimum number of coins needed to make the given amount.
// If the amount cannot be made with the given denominations, return -1.
func MinCoins(amount int, denominations []int) int {
	if amount == 0 {
		return 0
	}
	if amount < denominations[0] {
		return -1
	}
	count := 0
	for _, denomination := range reverse(denominations) {
		if amount >= denomination {
			count += amount / denomination
			amount %= denomination
		}
	}
	return count
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	result := make(map[int]int)
	for _, denomination := range reverse(denominations) {
		if amount >= denomination {
			result[denomination] = amount / denomination
			amount %= denomination
		}
	}
	return result
}
