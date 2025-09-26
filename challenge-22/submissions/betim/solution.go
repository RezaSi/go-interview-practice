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
		// Find coin combination
		coinCombo := CoinCombination(amount, denominations)

		minCoins := MinCoins(amount, denominations)

		// Print results
		fmt.Printf("Amount: %d cents\n", amount)
		fmt.Printf("Minimum coins needed: %d\n", minCoins)
		fmt.Printf("Coin combination: %v\n", coinCombo)
		fmt.Println("---------------------------")
	}
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) (combinations map[int]int) {
	combinations = make(map[int]int)

	for i := len(denominations) - 1; i >= 0; i-- {
		coin := denominations[i]
		count := amount / denominations[i]
		if count == 0 {
			continue
		}

		// fmt.Printf("%d / %d = %d\n", amount, denominations[i], compute)
		combinations[coin] = count
		amount = amount % coin

	}

	if amount != 0 {
		return make(map[int]int)
	}

	return
}

func MinCoins(amount int, denominations []int) (minCoins int) {
	if amount == 0 {
		return
	}

	combos := CoinCombination(amount, denominations)
	if len(combos) == 0 {
		return -1
	}

	for _, count := range combos {
		minCoins += count
	}

	return
}
