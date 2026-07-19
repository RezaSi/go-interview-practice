package main

import (
	"errors"
	"fmt"
	"math"
)

var (
	errOverflow = errors.New("errOverflow")
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
	result, err := Sum(a, b)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}

// Sum returns the sum of a and b.
func Sum(a int, b int) (int, error) {
	if a > 0 && b > math.MaxInt-a || a < 0 && b < math.MinInt-a {
		return 0, errOverflow
	}
	return a + b, nil
}
