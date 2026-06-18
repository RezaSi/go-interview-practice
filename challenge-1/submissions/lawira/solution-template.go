package main

import (
	"fmt"
)

func main() {
	var a, b int64
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
func Sum(a, b int64) int64 {
	// TODO: Implement the function
	result := a + b
	return result
}
