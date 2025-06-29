package main

import (
	"bufio"
	"fmt"
	"os"
	"bytes"
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
	// TODO: Implement the function
	buf := []byte(s)
	var reversed bytes.Buffer

	for len(buf) > 0 {
		r, size := utf8.DecodeLastRune(buf)
		reversed.WriteRune(r)
		buf = buf[:len(buf)-size]
	}
	
	return reversed.String()
}
