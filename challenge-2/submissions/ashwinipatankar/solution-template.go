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
	r := ""
	splittedString := strings.Split(s, "")

	le := len(splittedString) - 1
	for _ = range splittedString {
		r = r + splittedString[le]
		le--
	}

	return r
    
}
