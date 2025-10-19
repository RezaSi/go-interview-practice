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
	// 1. Clean the string (remove spaces, punctuation, and convert to lowercase)
	s = strings.ToLower(s)
	runes := []rune(s)
	res := make([]rune, 0, len(s))
	for _, r := range runes {
		if unicode.In(r, unicode.Letter, unicode.Digit) {
			res = append(res, r)
		}
	}
	// 2. Check if the cleaned string is the same forwards and backwards
	l := len(res)
	for i := 0; i < l/2; i++ {
		if res[i] != res[l-i-1] {
			return false
		}
	}

	return true
}
