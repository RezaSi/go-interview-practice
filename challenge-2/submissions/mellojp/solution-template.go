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
	runeSlice := []rune(s)
	size := len(runeSlice)
	res := make([]rune, size)
	i, j := 0, size-1
	for i < size && j >= 0 {
		res[i] = runeSlice[j]
		i++
		j--
	}
	return string(res)
}
