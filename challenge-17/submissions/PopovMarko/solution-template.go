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
	if s == "" {
		return true
	}

	// Normalization of the string to lower case, delete all not alphanumeric symbols
	// and spaces
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return ' '
	}, s)
	s = strings.ReplaceAll(s, " ", "")
	l := 0
	r := len(s) - 1

	// Loop with two indexes left and right that steps towards each other
	// and compare each symbol
	for l <= r {
		if s[l] != s[r] {
			return false
		}
		l++
		r--
	}
	return true
}
