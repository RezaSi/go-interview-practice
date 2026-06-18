package main

import (
	"fmt"
	"strings"
	"unicode"
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
    
	var cleaned strings.Builder
	
	for _, r := range s {
		lower := unicode.ToLower(r)

		if unicode.IsLetter(lower) || unicode.IsDigit(lower) {
			cleaned.WriteRune(lower)
		}
	}

	cleanedStr := cleaned.String()
    
	n := len(cleanedStr)
	
	for i := 0; i < n/2; i++ {
		if cleanedStr[i] != cleanedStr[n-1-i] {
			return false 
		}
	}

	return true
}
