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
	b := []rune(s)
	len := len(b)
	half := len/2
	for i:=0;i<half;i++{
	    b[i],b[len-i-1] = b[len-i-1],b[i]
	}
	return string(b)
}
