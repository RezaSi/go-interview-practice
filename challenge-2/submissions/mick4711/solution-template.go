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
	// check len contraint
    if len(s) >1000 {
        return "input string exceeds length limit (1000)"
    }
    
	// convert to slice of runes
	r := []rune(s)
	l := len(r)
	
	// swap in place
	for i :=0; i<l/2; i++ {
	   r[i], r[l-i-1] = r[l-i-1], r[i]
	}
	
// 	// make a slice for the result
// 	res := make([]rune, l)
	
// 	// fill result from the end backwards
// 	for i, c := range r {
// 	    res[l-i-1] = c
// 	}
	
	return string(r)
}
