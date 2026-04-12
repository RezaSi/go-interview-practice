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
	input := []rune(s)
	length := len(input)
	result := make([]rune, length)

	for idx := range input {
		result[idx] = input[length-idx-1]
	}

	return string(result)
}
