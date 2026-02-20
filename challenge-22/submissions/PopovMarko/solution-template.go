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
	// Slice for storage index - sum value - number of couns
	nc := make([]int, amount+1)
	for i := range nc {
		nc[i] = amount + 1
	}
	nc[0] = 0
	// Loop for ammount variants from 1 to ammount
	for i := 1; i <= amount; i++ {
		// Loop for coins to define the best variants
		for _, coin := range denominations {
			if i >= coin && nc[i-coin]+1 < nc[i] {
				nc[i] = nc[i-coin] + 1
			}
		}
	}
	if nc[amount] > amount {
		return -1
	}
	return nc[amount]
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	// Slice for storage index - sum value - number of couns
	nc := make([]int, amount+1)
	for i := range nc {
		nc[i] = amount + 1
	}
	nc[0] = 0
	// Slice for storage index - sum value - last added coin
	lc := make([]int, amount+1)
	// Map for storage result key - coin value - number of coins
	res := make(map[int]int)

	// Loop for ammount variants from 1 to ammount
	for i := 1; i <= amount; i++ {
		// Loop for coins to define the best variants
		for _, coin := range denominations {
			if i >= coin && nc[i-coin]+1 < nc[i] {
				nc[i] = nc[i-coin] + 1
				lc[i] = coin
			}
		}
	}
	if nc[amount] > amount {
		return map[int]int{}
	}
	currentSum := amount
	for currentSum > 0 {
		coin := lc[currentSum]
		res[coin]++
		currentSum -= coin
	}

	return res
}
