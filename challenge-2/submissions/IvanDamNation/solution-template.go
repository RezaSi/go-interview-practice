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
	
	for left, right := 0, len(runes)-1; left < right; left, right = left+1, right-1 {
	    runes[left], runes[right] = runes[right], runes[left]
	}
	
	return string(runes)
}
