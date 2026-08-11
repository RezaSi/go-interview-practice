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
	ret := 0
	for i := range denominations {
	    c := denominations[len(denominations)-i-1]
	    ret += amount / c
	    amount = amount % c 
	}
	if amount > 0 {
	    return -1
	}
	return ret
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	ret := make(map[int]int, 0)
	for i := range denominations {
	    c := denominations[len(denominations)-i-1]
	    ac := amount / c
	    if ac == 0 {
	        continue
	    }
	    if retc, ok := ret[c]; ok {
	        ret[c] = retc + ac
	    } else {
	        ret[c] = ac
	    }
	    amount = amount % c
	}
	if amount > 0 {
	    return map[int]int{}
	}
	return ret
}
