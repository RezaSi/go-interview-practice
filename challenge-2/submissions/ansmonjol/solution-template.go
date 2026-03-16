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
	r := []rune(s)
	left := 0
	right := len(r) - 1

	for left < right {
		r[left], r[right] = r[right], r[left]
		left++
		right--
	}

	return string(r)
}
