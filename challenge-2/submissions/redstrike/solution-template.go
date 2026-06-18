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
    // Swap runes or two-pointer technique
	runes := []rune(s) // Convert string to slice of runes
	runesLen := len(runes)
	halfRunesLen := runesLen / 2
	
	for i := 0; i < halfRunesLen; i++ {
		r1 := runes[i]
		r2 := runes[runesLen-1-i]
		runes[i] = r2
		runes[runesLen-1-i] = r1
	}
	
	return string(runes)
}
