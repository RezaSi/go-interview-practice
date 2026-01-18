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
    reversed := []byte(s)
    
	var j = len(reversed) - 1
	for i := 0; i < j; i++ {
	    store := reversed[i]
	    reversed[i] = reversed[j]
	    reversed[j] = store
	    j--
	}
	
	return string(reversed)
}
