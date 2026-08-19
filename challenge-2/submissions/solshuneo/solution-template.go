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
    mybyte := []rune(s)
    n := len(mybyte)
    for l, r := 0, n - 1; l < r; l, r = l + 1, r - 1 {
        mybyte[l], mybyte[r] = mybyte[r], mybyte[l]
    }
    
	return string(mybyte)
}
