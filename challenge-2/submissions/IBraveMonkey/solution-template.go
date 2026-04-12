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
	left := 0
	right := len(runes) -1
	
	for left < right {
	    runes[left], runes[right] = runes[right], runes[left]
	    left++
	    right--
	}
	
	return string(runes)
}
