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
	coins := 0
	
	for i := len(denominations)-1; i >=0; i-- {
	    denomination := denominations[i]
	    coinsToAdd := amount / denomination
	    amount = amount - denomination * coinsToAdd
	    coins += coinsToAdd
	}
	
	if amount != 0 {
	    return -1
	}
	return coins
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	coinCombination := make(map[int]int, len(denominations))
	
	for i := len(denominations)-1; i >=0; i-- {
	    denomination := denominations[i]
	    if coinsToAdd := amount / denomination; coinsToAdd > 0 {
	        coinCombination[denomination] += coinsToAdd
	        amount = amount - denomination * coinsToAdd
	    }
	}
	
	if amount != 0 {
	    return map[int]int{}
	}
	return coinCombination
}
