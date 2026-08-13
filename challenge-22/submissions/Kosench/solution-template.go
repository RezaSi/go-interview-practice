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

func MinCoins(amount int, denominations []int) int {
	if amount < 0 {
		return -1
	}
	if amount == 0 {
		return 0
	}

	sorted := append([]int(nil), denominations...)
	sort.Ints(sorted)

	totalCoins := 0
	remaining := amount

	for i := len(sorted) - 1; i >= 0; i-- {
		coin := sorted[i]
		if coin > 0 && coin <= remaining {
			count := remaining / coin
			totalCoins += count
			remaining %= coin
		}
		if remaining == 0 {
			break
		}
	}

	if remaining > 0 {
		return -1
	}
	return totalCoins
}

// CoinCombination returns a map with the specific combination of coins.
func CoinCombination(amount int, denominations []int) map[int]int {
	combo := make(map[int]int)
	if amount < 0 {
		return combo
	}

	sorted := append([]int(nil), denominations...)
	sort.Ints(sorted)

	remaining := amount

	for i := len(sorted) - 1; i >= 0; i-- {
		coin := sorted[i]
		if coin > 0 && coin <= remaining {
			count := remaining / coin
			combo[coin] = count
			remaining %= coin
		}
		if remaining == 0 {
			break
		}
	}

	if remaining > 0 {
		return make(map[int]int)
	}
	return combo
}
