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
	l := len(s)
	str := []rune(s)
    for i := 0; i < l / 2; i++ {
        temp := str[i]
        str[i] = str[l - 1 - i]
        str[l - 1 - i] = temp
    }
    
    return string(str)
}
