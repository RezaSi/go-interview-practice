package main

import (
	"fmt"
	"strings"
	"regexp"
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

// IsPalindrome checks if a string is a palindrome.
// A palindrome reads the same backward as forward, ignoring case, spaces, and punctuation.
func IsPalindrome(s string) bool {
 
    reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
    s = reg.ReplaceAllString(s, "")
    s = strings.ToLower(s)

	for i := 0; i < len(s) / 2; i++ {
	    if s[i] != s[len(s) - i - 1] {
	        return false
	    }
	}
	
	return true
}
