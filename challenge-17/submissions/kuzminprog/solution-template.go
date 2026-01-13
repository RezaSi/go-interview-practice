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
	cleaned := strings.Map(func(r rune) rune {
		if ch := unicode.ToLower(r); unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			return ch
		}
		return -1
	}, s)

	runes := []rune(cleaned)
	for l, r := 0, len(runes)-1; l < r; l, r = l+1, r-1 {
		if runes[l] != runes[r] {
			return false
		}
	}
	return true
}
