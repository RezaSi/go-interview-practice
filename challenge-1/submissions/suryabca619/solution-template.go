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

	// Call the Sum function and print the result
	result := Sum(a, b)
	fmt.Println(result)
}

// Sum returns the sum of a and b.
func Sum(a int, b int) int {
	return a + b
}
