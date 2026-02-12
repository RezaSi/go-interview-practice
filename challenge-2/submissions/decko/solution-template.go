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
    var reversed_rune []rune
    for i := len(runes)-1; i >= 0; i-- {
        reversed_rune = append(reversed_rune, runes[i])
    }
	return string(reversed_rune)
}
