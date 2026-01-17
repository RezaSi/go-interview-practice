package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	length := len(s)
	
	var sb strings.Builder
	sb.Grow(length)
	
	for i := 0; i < length; i ++ {
	    sb.WriteByte(s[length - i - 1])
	}
	
	return sb.String()
}
