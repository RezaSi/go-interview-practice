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
// Handle the base case 
if amount == 0{
return 0 
}
// Copy and sort the denominations in descending order 
coins := make([]int,len(denominations))
copy(coins, denominations)
num := 0 
sort.Sort(sort.Reverse(sort.IntSlice(coins)))
for _, coin := range coins {
    if amount>= coin {
        num += amount/coin
        amount = amount%coin
    }
}
if amount > 0 {
    return -1 
}
return num
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
// Copy and sort the denominations 
coins := make([]int, len(denominations))
copy(coins, denominations)

sort.Sort(sort.Reverse(sort.IntSlice(coins)))
ans := make(map[int]int, len(denominations))
for _, coin := range coins {
    if amount>= coin {
        
        ans[coin] = amount/coin
        amount = amount%coin
        
    }
}
if amount > 0 {
    return make(map[int]int)
}
	return ans 
}
