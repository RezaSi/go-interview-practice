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
	reverse := ""
	
	if len(s) > 0 && len(s) <= 1000 {
    	for i := len(s)-1; i >= 0; i-- {
    	    reverse += string(s[i])
    	}
	}
	
	return reverse
}
