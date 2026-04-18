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

// IsPalindrome checks if a string is a palindrome.
// A palindrome reads the same backward as forward, ignoring case, spaces, and punctuation.
func IsPalindrome(s string) bool {
	// TODO: Implement this function
	// 1. Clean the string (remove spaces, punctuation, and convert to lowercase)
	// 2. Check if the cleaned string is the same forwards and backwards
	s = strings.ToLower(s)
	var clean_string []byte
	for _, ch := range []byte(s) {
	    if IsAlphaNumeric(ch) {
	        clean_string = append(clean_string, ch)
	    }
	}
	
	i := 0
	j := len(clean_string) - 1
	for {
	    if i > j {
	        break
	    }
	    if clean_string[i] != clean_string[j] {
	        return false
	    }
	    i += 1
	    j -= 1
	}
	return true
}

func IsAlphaNumeric(ch byte) bool {
    return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')
}
