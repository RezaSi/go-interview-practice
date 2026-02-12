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
    text := []byte(s)
    textLength := len(s)
    var result string
    for i:=(textLength); i>0; i--{
        result += string(text[i-1])
    }
	// TODO: Implement the function
	return result
}
