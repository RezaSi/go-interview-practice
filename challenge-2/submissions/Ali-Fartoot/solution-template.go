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
	for startChar, endChar := 0, len(runes)-1; startChar < endChar; startChar, endChar = startChar+1, endChar-1 {
		runes[startChar], runes[endChar] = runes[endChar], runes[startChar]
	}
	
	result := string(runes)
	return result
}
