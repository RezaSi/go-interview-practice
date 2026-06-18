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
    sr := []rune(s)
	for i := 0; i < len(sr) >> 1; i++ {
	    sr[i], sr[len(sr) - 1 - i] = sr[len(sr) - 1 - i], sr[i]
	}
	return string(sr)
}
