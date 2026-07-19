package main

import (
  "fmt"
)

func main() {
  // 硬币面额数组
  denominations := []int{1, 5, 10, 25, 50}

  // 测试的金额数组
  amounts := []int{87, 42, 99, 33, 7}

  for _, amount := range amounts {
    // 找到满足金额使用最少的硬币数
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

// MinCoins 返回在给定的硬币面额数组denominations中，最少使用多少个硬币恰好等于amount金额
func MinCoins(amount int, denominations []int) int {
  cnt := 0
  for i := len(denominations) - 1; i >= 0; {
    if amount == 0 {
      return cnt
    }
    if amount < denominations[i] {
      i--
    } else {
      amount -= denominations[i]
      cnt++
    }
  }
  return -1
}

// CoinCombination 返回在给定的硬币面额数组denominations中，
// 每个面额最少使用多少个硬币使总额恰好等于amount金额
// 返回每种硬币的使用数量的map
func CoinCombination(amount int, denominations []int) map[int]int {
  cntMap := make(map[int]int)

  for i := len(denominations) - 1; i >= 0; {
    if amount == 0 {
      return cntMap
    }
    if amount < denominations[i] {
      i--
    } else {
      amount -= denominations[i]
      cntMap[denominations[i]]++
    }
  }
  return cntMap
}
