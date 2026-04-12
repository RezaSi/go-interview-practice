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

// MinCoins returns the minimum number of coins needed to make the given amount.
// If the amount cannot be made with the given denominations, return -1.
func MinCoins(amount int, denominations []int) int {
    // Create a copy of denominations and sort in descending order
    denoms := make([]int, len(denominations))
    copy(denoms, denominations)
    sort.Sort(sort.Reverse(sort.IntSlice(denoms)))
    
    coins := 0
    remaining := amount
    
    for _, coin := range denoms {
        if remaining <= 0 {
            break
        }
        if coin <= remaining {
            // Calculate how many of this coin we can use
            count := remaining / coin
            coins += count
            remaining -= count * coin
        }
    }
    
    // If we couldn't make the exact amount, return -1
    if remaining != 0 {
        return -1
    }
    
    return coins
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
    // Create a copy of denominations and sort in descending order
    denoms := make([]int, len(denominations))
    copy(denoms, denominations)
    sort.Sort(sort.Reverse(sort.IntSlice(denoms)))
    
    coinMap := make(map[int]int)
    remaining := amount
    
    for _, coin := range denoms {
        if remaining <= 0 {
            break
        }
        if coin <= remaining {
            // Calculate how many of this coin we can use
            count := remaining / coin
            coinMap[coin] = count
            remaining -= count * coin
        }
    }
    
    // If we couldn't make the exact amount, return empty map
    if remaining != 0 {
        return make(map[int]int)
    }
    
    return coinMap
}