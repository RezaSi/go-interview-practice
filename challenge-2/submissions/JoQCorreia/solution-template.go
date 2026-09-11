package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"slices"
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
	
	str := strings.Split(s, "")
	
	slices.Reverse(str)
	
	final := strings.Join(str, "")
	
	return final
}
