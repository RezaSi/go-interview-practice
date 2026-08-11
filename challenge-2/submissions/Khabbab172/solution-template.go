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
    r := []rune(s)
    j := len(r)-1
    i := 0
    for i < j {
        r[i] , r[j] = r[j] , r[i]
        i++
        j--
    } 
	return string(r)
}
