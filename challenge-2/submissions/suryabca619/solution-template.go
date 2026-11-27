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
    strRune := []rune(s)
	strLength := len(strRune)
	
	for i := 0; i < strLength / 2; i++ {
	    strRune[i], strRune[strLength - i - 1] = strRune[strLength - i - 1], strRune[i]
	}
	return string(strRune)
}
