package main

import (
	"bufio"
	"fmt"
	"strings"
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
    arr := strings.Split(s, "")
    var reverse string
    
    for i := len(arr)-1; i>=0;i--{
        reverse += arr[i]
    }
    
	return reverse
}
