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
		count := len(s)
	newString := make(map[int]string)
	for index, value := range s {
		newString[count-index-1] = string(value)
	}
	string2 := ""
	for i := 0; i < count; i++ {
		string2 += newString[i]
	}
	return string2
}
