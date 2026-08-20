package main

import (
	"fmt"
	"strings"
)

func main() {
	// Get input from the user
	var input string
	fmt.Print("Enter a string to check if it's a palindrome: ")
	fmt.Scanln(&input)

	// Call the IsPalindrome function and print the result
	result := IsPalindrome(input)
	if result {
		fmt.Println("The string is a palindrome.")
	} else {
		fmt.Println("The string is not a palindrome.")
	}
}


func IsPalindrome(s string) bool {

	cleanedstring := strings.ToLower(cleanstring(s))
	for i:=0;i<len(cleanedstring)/2;i++{
	    if cleanedstring[i] != cleanedstring[len(cleanedstring)-1-i] {
	        return false
	    }
	}
	return true
}

func cleanstring(s string) string {
    var result strings.Builder
    for _, r := range s {
        if ( r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
            result.WriteRune(r)
        }
    
    }
    return result.String()
    
}
