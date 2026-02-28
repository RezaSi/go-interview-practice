package main

import (
	"errors"
	"fmt"
	"math"
	"os"
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
		os.Exit(1)
	}

	// Call the Sum function and print the result
	result, err := Sum(a, b)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
