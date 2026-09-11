package main

import (
	"fmt"
	"os"
)

func main() {
	var a, b int
	_, err := fmt.Scanf("%d, %d", &a, &b)
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}
	result := Sum(a, b)
	fmt.Println(result)
}

func Sum(a, b int) int {
	return a + b
}