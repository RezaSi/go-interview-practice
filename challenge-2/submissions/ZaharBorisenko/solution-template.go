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
	if s == "" {
		return s
	}
	result := ""
	data := []byte(s)

	for i := len(data) - 1; i >= 0; i-- {
		result += string(data[i])
	}
	return result
}

