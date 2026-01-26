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
	j := len(s) - 1
	for i := 0; i< len(s)/2; i++{
	    runes[i], runes[j] = runes[j], runes[i]
	    j--
	}
	return string(runes)
}
