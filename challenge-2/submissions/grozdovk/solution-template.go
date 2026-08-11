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
    left,right:= 0,len(s)-1
    if len(s)==-1{
	return ""
    }
    runes := []rune(s)
    result:= make([]rune, right+1)
    for left<=right{
        result[left] =runes[right]
        result[right] = runes[left]
        left++
        right--
    }
    
    return string(result)
}
