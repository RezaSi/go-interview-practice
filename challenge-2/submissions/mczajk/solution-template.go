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
    if len(s) == 0 {
        return ""
    }  else {
    	runes := []rune(s)    // Convert string to slice of runes
    
        resultRunes := []rune("")
    
        for i := len(runes)-1; i >=0 ; i-- {
            resultRunes = append(resultRunes, runes[i])
        }
        return string(resultRunes) 
    }
}
