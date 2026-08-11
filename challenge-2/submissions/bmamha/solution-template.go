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
	rune_s := []rune(s)
	left := 0
	right := len(rune_s) - 1
	for left < right {
	    rune_s[left], rune_s[right] = rune_s[right], rune_s[left]
	    left++
	    right--
	}
	return string(rune_s) 
}
