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

// getDenominationUsage returns the map that have information about coins needed to make the given amount.
func getDenominationUsage(amount int, denominations []int)map[int]int{
    sort.Slice(denominations, func(i, j int) bool { return denominations[i] > denominations[j] })
    usedDenominations := make(map[int]int)
    for _, d := range denominations {
		coinToUse := amount / d
		if coinToUse > 0 {
			usedDenominations[d] = coinToUse
			amount = amount % d
		}
	}
	return usedDenominations
}

// MinCoins returns the minimum number of coins needed to make the given amount.
// If the amount cannot be made with the given denominations, return -1.
func MinCoins(amount int, denominations []int) int {
	// TODO: Implement this function
	if amount == 0{
	    return 0
	}
	usedDenominations := getDenominationUsage(amount,denominations)
	if usedDenominations == nil || len(usedDenominations)==0{
	    return -1
	}
	total := 0;
	for _,v := range usedDenominations{
	    total=total+v
	}
	return total
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	usedDenominations := getDenominationUsage(amount,denominations)
	return usedDenominations
}
