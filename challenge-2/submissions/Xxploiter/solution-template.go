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
// 	first i will try to convert the string to rune
 cnvStr:= []rune(s)
 strLen:= len(cnvStr)
 lasteEle:= strLen - 1
 for i:= 0; i<= (strLen/2)-1; i++ {
     cnvStr[i] ,cnvStr[lasteEle-i] = cnvStr[lasteEle-i] ,cnvStr[i]
 }
 strString := string(cnvStr)
 
	return strString
}
