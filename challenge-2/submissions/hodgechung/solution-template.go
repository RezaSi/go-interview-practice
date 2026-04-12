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
    if (len(s) == 0) {
        return s
    }
    ru := []rune(s)
    for l, r := 0, len(ru)-1; l < r;  {
        ru[l], ru[r] = ru[r], ru[l]
        l++
        r--
    }
	return string(ru)
}
