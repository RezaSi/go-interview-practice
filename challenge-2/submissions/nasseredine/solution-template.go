package main

import (
	"bufio"
	"fmt"
	"os"
	
	"golang.org/x/example/hello/reverse"
)

// 
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
    for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
	return reverse.String(s)
}