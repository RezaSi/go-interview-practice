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
	size := len(runes)
	
	for i:=0; i < size / 2; i++ {
	    runes[i], runes[size-i-1] = runes[size-i-1], runes[i]
	}
	
	return string(runes)
}
