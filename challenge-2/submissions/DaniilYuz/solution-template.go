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
	res := []byte(s)
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
	    res[i], res[j] = res[j], res[i] 
	}
	return string(res)
}
