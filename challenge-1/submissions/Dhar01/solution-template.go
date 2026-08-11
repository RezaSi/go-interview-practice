package main

import (
	"fmt"
	"math"
)

func main() {
	var a, b int
	// Read two integers from standard input
	_, err := fmt.Scanf("%d, %d", &a, &b)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Call the Sum function and print the result
	result := Sum(a, b)
	fmt.Println(result)
}

// Sum returns the sum of a and b.
func Sum(a int, b int) int {
	// TODO: Implement the function
	lowerLimit := math.Pow(-10, 9)
	higherLimit := math.Pow(10, 9)
	
	if int(lowerLimit) <= a && b <= int(higherLimit) {
	    return a + b
	}
	
	return 0
}
