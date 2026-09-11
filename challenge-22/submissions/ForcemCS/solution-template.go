package main

import (
	"fmt"
)

func main() {

	denominations := []int{1, 5, 10, 25, 50}

	amounts := []int{87, 42, 99, 33, 7}

	for _, amount := range amounts {

		minCoins := MinCoins(amount, denominations)

		coinCombo := CoinCombination(amount, denominations)

		fmt.Printf("Amount: %d cents\n", amount)
		fmt.Printf("Minimum coins needed: %d\n", minCoins)
		fmt.Printf("Coin combination: %v\n", coinCombo)
		fmt.Println("---------------------------")
	}
}

func MinCoins(amount int, denominations []int) int {
	totalCoins := 0

	for i := len(denominations) - 1; i >= 0; i-- {
		coin := denominations[i]

		if amount >= coin {
			count := amount / coin
			totalCoins += count
			amount = amount % coin
		}
	}

	if amount > 0 {
		return -1
	}

	return totalCoins
}

func CoinCombination(amount int, denominations []int) map[int]int {
	combination := make(map[int]int)

	for i := len(denominations) - 1; i >= 0; i-- {
		coin := denominations[i]

		if amount >= coin {
			count := amount / coin
			combination[coin] = count
			amount = amount % coin

		}

	}

	if amount > 0 {
		return map[int]int{}
	}

	return combination
}
