package main

import (
	"fmt"
)

func main() {
	// Standard U.S. coin denominations in cents
	// 题目要求：切片是按升序排列的
	denominations := []int{1, 5, 10, 25, 50}

	// Test amounts from Sample Input
	amounts := []int{87, 42}

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
	
	// 测试无法凑出的情况 (例如没有1分硬币，只要凑3分)
	// unmatchableCoins := MinCoins(3, []int{2, 5})
	// fmt.Printf("Test Unmatchable (3 with [2,5]): %d\n", unmatchableCoins)
}

// MinCoins returns the minimum number of coins needed to make the given amount.
// If the amount cannot be made with the given denominations, return -1.
func MinCoins(amount int, denominations []int) int {
	count := 0
	remaining := amount

	// 贪心算法：从最大面额（切片末尾）开始向前遍历
	for i := len(denominations) - 1; i >= 0; i-- {
		coin := denominations[i]
		
		// 只要剩余金额大于等于当前硬币面额
		if remaining >= coin {
			// 使用除法计算当前面额需要几枚
			num := remaining / coin
			
			// 累加硬币数量
			count += num
			
			// 更新剩余金额 (取余数)
			remaining %= coin
		}
	}

	// 题目要求：如果无法凑出该金额，返回 -1
	if remaining != 0 {
		return -1
	}

	return count
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
	// 初始化 map
	sol := make(map[int]int)
	remaining := amount

	// 贪心算法：从最大面额开始
	for i := len(denominations) - 1; i >= 0; i-- {
		coin := denominations[i]
		
		if remaining >= coin {
			num := remaining / coin
			
			if num > 0 {
				sol[coin] = num
				remaining %= coin
			}
		}
	}

	// 题目要求：如果无法凑出，返回空 map
	if remaining != 0 {
		return map[int]int{}
	}

	return sol
}