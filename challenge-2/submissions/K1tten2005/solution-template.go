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
    runeStr := []rune(s)
    l, r := 0, len(runeStr) - 1
    for l < r {
        runeStr[l], runeStr[r] = runeStr[r], runeStr[l]
        l++
        r--
    }
    return string(runeStr)
}
