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
	runes := []rune(s)
	l := len(runes)
	if l <= 1 {
	    return s
	} 
    for i := 0; i < l / 2; i++ {
        temp := runes[i]
        runes[i] = runes[l - 1 - i]
        runes[l - 1 - i] = temp
    }
    result := string(runes)
	return result
}
