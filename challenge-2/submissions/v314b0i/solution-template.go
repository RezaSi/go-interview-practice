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
	chars := []rune(s)
	length := len(chars)
	if length <= 1 {
	    return s
	}
	for i := 0; i< length/2; i+=1{
	    chars[i], chars[length-i-1] = chars[length-i-1], chars[i]
	}
	return string(chars)
}
