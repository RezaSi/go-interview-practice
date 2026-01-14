package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode/utf8"
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
	n := len(s)
	buf := make([]byte, n)
	start := 0

	for start < n {
		_, size := utf8.DecodeRuneInString(s[start:])
		copy(buf[n-start-size:n-start], s[start:start+size])

		start += size
	}

	return string(buf)
}
