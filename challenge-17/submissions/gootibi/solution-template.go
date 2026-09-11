package main

import (
	"fmt"
	"unicode"
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
	// 1. Clean the string (remove spaces, punctuation, and convert to lowercase)
	builder := strings.Builder{}

	for _, str := range s {
		if unicode.IsLetter(str) || unicode.IsDigit(str) {
			builder.WriteRune(unicode.ToLower(str))
		}
	}

	// 2. Check if the cleaned string is the same forwards and backwards
	cleanStr := builder.String()

	for i := 0; i < len(cleanStr)/2; i++ {
		if cleanStr[i] != cleanStr[len(cleanStr)-i-1] {
			return false
		}
	}

	return true
}
