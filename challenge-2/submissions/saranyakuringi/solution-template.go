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
	run:=[]rune(s)
	left,right:=0,len(run)-1
	for i:=0;i<len(run)/2;i++{
	    run[left],run[right]=run[right],run[left]
	    left++
	    right--
	}
	return string(run)
}
