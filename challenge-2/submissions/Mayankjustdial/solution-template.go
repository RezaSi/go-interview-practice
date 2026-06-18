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
	runes := []rune(s)
	start, end := 0, len(runes)-1

	for start < end {
		runes[start], runes[end] = runes[end], runes[start]
		start++
		end--
	}

	return string(runes)
}
