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
	result := ""
	runeSlice := []rune(s)

	for i := len(runeSlice) - 1; i >= 0; i-- {
		result += fmt.Sprintf("%c", runeSlice[i])
	}
	return result
}
