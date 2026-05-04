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
func ReverseString(str string) string {

    req := []rune(str);
    
    for i:=0;i<len(str)/2;i++{
        req[i],req[len(str)-i-1]=req[len(str)-i-1],req[i];
    }
    
	
	return string(req);
}
