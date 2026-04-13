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

// initHelper returns two slices nc - with the best numbers of coins for ech sum
// and lc - with the last coin added for sum
func initHelper(amount int, denominations []int) ([]int, []int) {
	// Slice for storage index - sum value - number of coins
	nc := make([]int, amount+1)
	// Init slice by max value
	for i := range nc {
		nc[i] = amount + 1
	}
	// It meens for 0 sum 0 coins
	nc[0] = 0
	// Slice for storage index - sum value - last added coin
	lc := make([]int, amount+1)

	// Loop for amount variants from 1 to amount
	for i := 1; i <= amount; i++ {
		// Loop for coins to define the best variants
		for _, coin := range denominations {
			if i >= coin && nc[i-coin]+1 < nc[i] {
				nc[i] = nc[i-coin] + 1
				lc[i] = coin
			}
		}
	}
	return nc, lc
}

// MinCoins returns the minimum number of coins needed to make the given amount.
// If the amount cannot be made with the given denominations, return -1.
func MinCoins(amount int, denominations []int) int {
	if amount == 0 {
		return 0
	}
	nc, _ := initHelper(amount, denominations)
	// Check if unpossible to make a given amount returns int
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
	if amount == 0 {
		return map[int]int{}
	}
	// Map for storage result key - coin value - number of coins
	res := make(map[int]int)
	nc, lc := initHelper(amount, denominations)
	// Check if unpossible to make given amount returns map[int]int
	if nc[amount] > amount {
		return map[int]int{}
	}
	// Loop for make result map from lc slice returns res - map[int]int
	currentSum := amount
	for currentSum > 0 {
		coin := lc[currentSum]
		res[coin]++
		currentSum -= coin
	}

	return res
}
