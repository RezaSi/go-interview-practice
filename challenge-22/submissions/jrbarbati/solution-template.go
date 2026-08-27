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

func MinCoins(amount int, denominations []int) int {
    var sum int
    
    for i := len(denominations)-1; i >= 0; i-- {
        sum += amount / denominations[i]
        amount %= denominations[i]
    }
    
    if amount > 0 {
        return -1
    }
    
    return sum
}

func CoinCombination(amount int, denominations []int) map[int]int {
    combo := make(map[int]int)
    
    for i := len(denominations)-1; i >= 0; i-- {
        if amount / denominations[i] > 0 {
            combo[denominations[i]] = amount / denominations[i]
            amount %= denominations[i]
        }
    }
    
    if amount > 0 {
        return map[int]int{}
    }
    
    return combo
}
