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
    oldString := []rune(s)
	newString := make([]rune, len(oldString))
	stringLen := len(oldString)
	for i := 0; i < stringLen; i++ {
		changI := stringLen - (i + 1)
		newString[i] = oldString[changI]
	}
	return string(newString)
}
