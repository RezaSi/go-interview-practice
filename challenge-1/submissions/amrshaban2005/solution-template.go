package main

import (
	"fmt"
)

func main() {
	var a, b int
	// Read two integers from standard input
	_, err := fmt.Scanf("%d, %d", &a, &b)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	const min = -1000_000_000
	const max = 1000_000_000
	if a < min || a > max || b < min || b > max {
		fmt.Println("invalid inputs")
		return
	}

	// Call the Sum function and print the result
	result := Sum(a, b)
	fmt.Println(result)
}

// Sum returns the sum of a and b.
func Sum(a int, b int) int {
	return a + b
}
