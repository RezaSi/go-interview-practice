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
	ret := make([]byte, len(s))
	i := 0
	j := len(s) - 1
	for i <= j {
		ret[i] = s[j]
		ret[j] = s[i]
		i++
		j--
	}
	return string(ret)
}
