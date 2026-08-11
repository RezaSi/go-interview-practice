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
	sArr := []rune(s)
	i, j := 0, len(s)-1
	for i < j {
		sArr[i], sArr[j] = sArr[j], sArr[i]
		i++
		j--
	}
	return string(sArr)
}
