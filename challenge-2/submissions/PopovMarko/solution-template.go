package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Read input from standard input
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading input:", err)
			os.Exit(1)
		}
	}

	input := scanner.Text()

	// Call the ReverseString function
	output := ReverseString(input)

	// Print the result
	fmt.Println(output)
}

// ReverseString returns the reversed string of s.
func ReverseString(s string) string {
	runes := []rune(s)
	// Loop runes with two indexes to reverse on place
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
