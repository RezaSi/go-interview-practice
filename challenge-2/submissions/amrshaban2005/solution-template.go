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
		if len(input) == 0 || len(input) > 1000 {
			fmt.Println("invalid string length")
			return
		}
		// Call the ReverseString function
		output := ReverseString(input)

		// Print the result
		fmt.Println(output)
	}
}

// ReverseString returns the reversed string of s.
func ReverseString(s string) string {

	text := []rune(s)

	var reversed string

	for i := len(s) - 1; i >= 0; i-- {
		reversed += string(text[i])
	}

	return reversed
}
