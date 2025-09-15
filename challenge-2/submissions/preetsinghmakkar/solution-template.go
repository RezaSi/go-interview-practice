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
    
    if(len(s) == 0){
        return ""
    }

	n := len(s)

	result := ""

	for i := n - 1; i >= 0; i-- {
		result += string(s[i])
	}

	fmt.Println(result)
    
	return result
}
