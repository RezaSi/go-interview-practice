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
	runes := []rune(s)
	i := 0
	j := len(runes)-1
	for i < j {
	    runes[i],runes[j] = runes[j],runes[i]
	    i++
	    j--
	}
	return string(runes)
}
