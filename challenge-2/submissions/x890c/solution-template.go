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
    t := []rune(s)
	r := make([]rune, len(t))
	for i:=0;i<len(t);i++ {
	    r[i] = t[len(t) - i -1]
	}
	return string(r)
}
