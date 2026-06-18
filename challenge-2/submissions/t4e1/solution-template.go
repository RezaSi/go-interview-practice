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
	
	for i := 0; i < len(runes)/2; i++ {
	    x := len(runes) - (i+1) 
	    runes[i], runes[x] = runes[x], runes[i]
	}
	
	return string(runes)
}
