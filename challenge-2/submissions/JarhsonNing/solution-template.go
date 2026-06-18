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
	runeSlice := []rune(s)
	resultSlice := []rune(s)
	for i := len(runeSlice) - 1; i >= 0; i-- {
		resultSlice[len(runeSlice)-1-i] = runeSlice[i]
	}
	return string(resultSlice)
}
