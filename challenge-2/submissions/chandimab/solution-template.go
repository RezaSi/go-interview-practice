package main

import (
	"bufio"
	"fmt"
	"os"
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

// ReverseString returns the reversed string of s.
func ReverseString(s string) string {
    // convert the string to rune to get positional character values.
    //runes := [] rune(s) // array of runes
    var sb strings.Builder // let's use string builder
    
    for i:= len(s) - 1; i >= 0; i-- {
        //sb.WriteString(string(runes[i])) using runes
        sb.WriteString(s[i: i + 1])
    }
	return sb.String()
}
