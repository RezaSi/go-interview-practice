package main

import "fmt"

func main() {
	var a, b int

	_, err := fmt.Scanf("%d,%d", &a, &b)
	if err != nil {
		fmt.Println("Error reading input", err)
		return
	}

	result := Sum(a, b)
	fmt.Println(result)
}

// Sum returns the sum of a and b.
func Sum(a int, b int) int {
	return a + b
}