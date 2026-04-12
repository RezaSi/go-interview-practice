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
	
	start := 0
	end := len(s)-1
	
	r := []rune(s)
	
	for start <= end {
	    r[start], r[end] = r[end], r[start]
	    start++
	    end--
	}
	return string(r)
}
