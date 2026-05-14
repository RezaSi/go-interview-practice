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
	s1 := []rune(s)
	for i,j:= 0,len(s1)-1;i<j;i,j=i+1,j-1{
	    s1[i],s1[j]=s1[j],s1[i]
	}
	return string(s1)
}
