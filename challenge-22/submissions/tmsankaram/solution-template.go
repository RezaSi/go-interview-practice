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
	minCoins := make([]int, amount+1)

	for i := range minCoins {
		minCoins[i] = -1
	}

	minCoins[0] = 0

	for _, coin := range denominations {
		for j := coin; j <= amount; j++ {
			if minCoins[j-coin] != -1 {
				if minCoins[j] == -1 || minCoins[j] > minCoins[j-coin]+1 {
					minCoins[j] = minCoins[j-coin] + 1
				}
			}
		}
	}

	if minCoins[amount] != -1 {
		return minCoins[amount]
	}

	return -1
}

func CoinCombination(amount int, denominations []int) map[int]int {

	coinCombo := make(map[int]int)

	if MinCoins(amount, denominations) == -1 {
		return coinCombo
	}

	for amount > 0 {
		for _, coin := range denominations {
			if amount-coin >= 0 && MinCoins(amount-coin, denominations) == MinCoins(amount, denominations)-1 {
				coinCombo[coin]++
				amount -= coin
				break
			}
		}
	}
	return coinCombo
}
