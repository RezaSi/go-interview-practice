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


func biggest(amount int, d []int) int{
    i:=len(d)-1
    if (amount-d[0])<0{
        return -1 
    }
    for ;i>0;i--{
        if(amount>=d[i]){
        return d[i]}
    }
    return d[i]
}
// MinCoins returns the minimum number of coins needed to make the given amount.
// If the amount cannot be made with the given denominations, return -1.
func MinCoins(amount int, denominations []int) int {
	// TODO: Implement this function
	if len(denominations)==0||amount==0{
	    return 0
	}
	if biggest(amount, denominations)==-1{
	    return -1
	}
	m:=0
	
	for {
	    b:=biggest(amount, denominations)
	    if b<denominations[0]||b<0{
	        break
	    }
	    amount=amount-biggest(amount, denominations)
	    m++
	}
	return m
}

// CoinCombination returns a map with the specific combination of coins that gives
// the minimum number. The keys are coin denominations and values are the number of
// coins used for each denomination.
// If the amount cannot be made with the given denominations, return an empty map.
func CoinCombination(amount int, denominations []int) map[int]int {
    c:=map[int]int{}
	for {
	    
	    b:=biggest(amount, denominations)
	    if b<0{
	        break
	    } else
	    if b<=amount{
	    c[b]++
	    amount-=b}else{
	        break
	    }
	    
	}
	return c
}
