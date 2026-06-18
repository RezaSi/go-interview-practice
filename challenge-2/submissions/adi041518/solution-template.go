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
	runes:=[]rune(s)
	low:=0
	high:=len(s)-1
	for low<=high{
	    var temp rune
	    temp=runes[low]
	    runes[low]=runes[high]
	    runes[high]=temp
        low++
        high--
	}
	return string(runes)
}
