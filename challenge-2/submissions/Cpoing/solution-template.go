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
	runes := []rune(s)

	for i := range runes {
		if i == len(s)/2 {
			break
		}

		tmp := runes[i]
		runes[i] = runes[len(s)-1-i]
		runes[len(s)-1-i] = tmp
	}
	return string(runes)
}
