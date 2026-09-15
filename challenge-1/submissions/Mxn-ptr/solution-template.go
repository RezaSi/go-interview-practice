package main

import (
	"fmt"
)

// main reads two integers from standard input,
// computes their sum using the Sum function,
// and prints the result.
func main() {
	var a, b int

	// Read two integers from standard input (format: "a, b")
	_, err := fmt.Scanf("%d, %d", &a, &b)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Call the Sum function and print the result
	result := Sum(a, b)
	fmt.Println(result)
}

// Sum returns the sum of two integers.
//
// Parameters:
//   - a: first integer
//   - b: second integer
//
// Returns:
//   - the sum of a and b
func Sum(a int, b int) int {
	return a + b
}
