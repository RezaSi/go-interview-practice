package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Read input from standard input
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input := scanner.Text()

		// Call the ReverseString function
		output := ReverseString(input)

		// Print the result
		fmt.Println(output)
	}
}

// ReverseString returns the reversed string of s.
func ReverseString(s string) string {
	// TODO: Implement the function
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s = swap(s, i, j)
	}

	return s
}

func swap(s string, i, j int) string {
	runes := []rune(s)
	runes[i], runes[j] = runes[j], runes[i]
	return string(runes)
}
