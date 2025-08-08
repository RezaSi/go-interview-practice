package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Item represents an item with weight and value
type Item struct {
	Weight int
	Value  int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Read number of items
	scanner.Scan()
	n, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Println("Error reading number of items:", err)
		return
	}

	// Read items
	items := make([]Item, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		parts := strings.Fields(scanner.Text())
		if len(parts) != 2 {
			fmt.Println("Error: expected weight and value")
			return
		}

		weight, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Println("Error reading weight:", err)
			return
		}

		value, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Error reading value:", err)
			return
		}

		items[i] = Item{Weight: weight, Value: value}
	}

	// Read capacity
	scanner.Scan()
	capacity, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		fmt.Println("Error reading capacity:", err)
		return
	}

	// Call the Knapsack function and print the result
	result := Knapsack(items, capacity)
	fmt.Println(result)
}

// Knapsack returns the maximum value that can be achieved
// without exceeding the given capacity.
func Knapsack(items []Item, capacity int) int {
	// TODO: Implement the function
	return 0
}
