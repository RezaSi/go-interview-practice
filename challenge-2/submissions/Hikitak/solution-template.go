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
	buf := strings.Builder{}
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
	    buf.WriteRune(runes[i])
	}
	return buf.String()
}
