package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
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

func reverseWord(s string) string {
	chars := strings.Split(s, "")
	start := 0
	end := len(s) - 1

	for start < len(s)/2 && start < end {
		tmp := chars[start]
		chars[start] = chars[end]
		chars[end] = tmp

		start++
		end--
	}
	return strings.Join(chars, "")
}

// ReverseString returns the reversed string of s.
func ReverseString(s string) string {
	// TODO: Implement the function

	// "Go is fun!"
	words := strings.Fields(s)
	reversedWorfs := make([]string, len(words))

	for _, w := range words {
		reversedWorfs = append(reversedWorfs, reverseWord(w))
	}

	slices.Reverse(reversedWorfs)

	return strings.Join(reversedWorfs, " ")
}
