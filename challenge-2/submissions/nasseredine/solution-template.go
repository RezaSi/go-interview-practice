package main

import (
	"bufio"
	"fmt"
	"os"
	
	"golang.org/x/example/hello/reverse"
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
    // Don't reinvent the wheel!
	return reverse.String(s)
}
