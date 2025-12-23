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
    b := []byte(s)
	for i := 0; i<(len(s)/2);i++{
	    c := b[i]
	    b[i] = b[len(b)-i-1]
	    b[len(b)-i-1] = c
	}
	return string(b)
}
