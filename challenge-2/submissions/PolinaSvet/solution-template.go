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

	data := []rune(s)
	lenData := len(data)
	result := make([]rune, lenData, lenData)

	for i, v := range data {
		result[lenData-i-1] = v
	}

	return string(result)
}
