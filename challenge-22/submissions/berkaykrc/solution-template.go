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
	if amount < 0 {
		return -1
	}
	coins := append([]int(nil), denominations...)
	sort.Sort(sort.Reverse(sort.IntSlice(coins)))
	remainingAmount := amount
	coinCount := 0
	for _, coin := range coins {
		count := remainingAmount / coin
		coinCount += count
		remainingAmount -= count * coin
	}
	if remainingAmount != 0 {
		return -1
	}
	return coinCount
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	if amount < 0 {
		return map[int]int{}
	}
	coins := append([]int(nil), denominations...)
	sort.Sort(sort.Reverse(sort.IntSlice(coins)))
	remainingAmount := amount
	combination := make(map[int]int)
	for _, coin := range coins {
		count := remainingAmount / coin
		if count > 0 {
			combination[coin] = count
		}
		remainingAmount -= count * coin
	}
	if remainingAmount != 0 {
		return map[int]int{}
	}
	return combination
}
