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
// ReverseString returns the reversed string of s.
func ReverseString(s string) string {

	reverse := []rune(s)
	length := len(reverse)

	for i := 0; i < length/2; i++ {
		reverse[i], reverse[length-i-1] = reverse[length-i-1], reverse[i]

	}
	return string(reverse)
}
