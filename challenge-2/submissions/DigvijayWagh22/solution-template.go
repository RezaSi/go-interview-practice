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
	runeArray := []rune(s)
	n := len(runeArray)
	low := 0
	high := n-1
	
	for low < high {
	    runeArray[low], runeArray[high] = runeArray[high], runeArray[low]
	    low++
	    high--
	}
	result := string(runeArray)
	return result
}
