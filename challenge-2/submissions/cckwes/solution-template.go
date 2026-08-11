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
    stringLength := len(runes)
    result := make([]rune, stringLength)
    
    for i := 0; i < stringLength; i++ {
        result[stringLength - 1 - i] = runes[i]
    }
    
	return string(result)
}
