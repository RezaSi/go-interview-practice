package main

import (
	"fmt"
	"strings"
	"regexp"
)

var (
    regex = regexp.MustCompile("[^A-Za-z0-9]")
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
    s = strings.ToLower(regex.ReplaceAllString(s, ""))
    var i = 0
    var j = len(s)-1
    
    for i < j {
        if s[i] != s[j] {
            return false
        }
        
        i++
        j--
    }
    
	return true
}
