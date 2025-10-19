package main

import (
	"fmt"
	"slices"
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
	// check for amount > 0
	if amount <= 0 {
		return 0
	}

	// init remainder and totalCount
	rem, totalCount := amount, 0

	// sort denominations in descending value order
	slices.SortFunc(denominations, func(a, b int) int {
		switch {
		case a > b:
			return -1
		case a < b:
			return 1
		}
		return 0
	})

	// iterate through denoms
	for _, denom := range denominations {
		// skip denom if too large
		if denom > rem {
			continue
		}

		// get number of denoms that can be used and accumulate
		count := rem / denom
		totalCount += count

		// decrease remainder by value taken
		rem -= denom * count

		// finish if none left
		if rem == 0 {
			return totalCount
		}
	}

	// couldn't use denoms for this amount
	return -1
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	// alocate results map
	res := make(map[int]int, len(denominations))

	// check for amount > 0
	if amount <= 0 {
		return res
	}

	// init remainder
	rem := amount

	// sort denominations in descending value order
	slices.SortFunc(denominations, func(a, b int) int {
		switch {
		case a > b:
			return -1
		case a < b:
			return 1
		}
		return 0
	})

	// iterate through denoms
	for _, denom := range denominations {
		// skip denom if too large
		if denom > rem {
			continue
		}

		// get number of denoms that can be used and assign to map
		count := rem / denom
		res[denom] = count

		// decrease remainder by value taken
		rem -= denom * count

		// finish if none left
		if rem == 0 {
			return res
		}
	}

	// couldn't use denoms for this amount, return empty map
	return res
}
