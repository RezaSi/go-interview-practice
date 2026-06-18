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
	arr := []rune(s)
	for left, right := 0, len(arr)-1; left < len(arr)/2; left, right = left+1, right-1 {
		arr[left], arr[right] = arr[right], arr[left]
	}
	return string(arr)
}
