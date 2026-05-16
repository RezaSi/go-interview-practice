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
	// TODO: Implement the function
	new_string := []int32(s)
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		new_string[i], new_string[j] = new_string[j], new_string[i]
	}
	return string(new_string)
}
