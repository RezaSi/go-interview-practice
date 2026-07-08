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
	if amount < 0 {
		return -1
	}
	if amount == 0 {
		return 0
	}

	totalCoins := 0
	remaining := amount

	// Идем с конца, так как жадный алгоритм требует начинать с самой крупной монеты
	for i := len(denominations) - 1; i >= 0; i-- {
		coin := denominations[i]
		if coin > 0 && coin <= remaining {
			count := remaining / coin
			totalCoins += count
			remaining %= coin // Эквивалентно remaining -= count * coin, но чище
		}
		if remaining == 0 {
			break
		}
	}

	// Если после прохода всех монет остаток всё еще есть, набрать сумму невозможно
	if remaining > 0 {
		return -1
	}
	return totalCoins
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	combo := make(map[int]int)
	if amount < 0 {
		return combo
	}

	remaining := amount

	for i := len(denominations) - 1; i >= 0; i-- {
		coin := denominations[i]
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
		return make(map[int]int) // Возвращаем пустой map, как требуется в условии
	}
	return combo
}
