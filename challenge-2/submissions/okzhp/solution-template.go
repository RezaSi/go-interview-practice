package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
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
	length := len(runes)
	for i := 0; i < length/2; i++ {
		tmp := runes[i]
		runes[i] = runes[length-i-1]
		runes[length-i-1] = tmp
	}
	return string(runes)
}
