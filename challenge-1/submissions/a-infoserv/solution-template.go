package main

import (
	"fmt"
)
var a, b int
// Sum returns the sum of a and b.
func Sum(a int, b int) int {
	
	return  a + b
}



func main() {
	
	// Read two integers from standard input
	_, err := fmt.Scanf("%d, %d", &a, &b)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Call the Sum function and print the result
	result := Sum(1, 2)
	fmt.Println(result)
}


