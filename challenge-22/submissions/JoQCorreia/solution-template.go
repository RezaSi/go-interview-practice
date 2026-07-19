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
func MinCoins(amount int, denominations []int) int {
    orderedDenom := slices.Clone(denominations)
    slices.Reverse(orderedDenom)
	sumDenom := 0
	sumFinal := 0

	if amount == 0 {
		return sumDenom
	}

	for i := range orderedDenom {
		for sumFinal + orderedDenom[i] <= amount {
			sumFinal = sumFinal + orderedDenom[i]
			sumDenom++
		}
	}

	if sumFinal != amount {
		return -1
	}

	return sumDenom
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
func CoinCombination(amount int, denominations []int) map[int]int {
	
	orderedDenom := slices.Clone(denominations)
	slices.Reverse(orderedDenom)
	
	denomMap := make(map[int]int)
	
	sumFinal := 0

	if amount == 0 {
		return denomMap
	}
	
	

	for i := range orderedDenom {
		for sumFinal + orderedDenom[i] <= amount {
			sumFinal = sumFinal + orderedDenom[i]
			denomMap[orderedDenom[i]] ++
		}
	}

	if sumFinal != amount {
		clear(denomMap)
		return denomMap
	}

	return denomMap
}
