package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
    var b strings.Builder
    r := []rune(s)
    l := len(r)
    b.Grow(l)
    
	for i:=l-1; i>=0;i-- {
	    b.WriteRune(r[i])   
	}
	return b.String()
}
