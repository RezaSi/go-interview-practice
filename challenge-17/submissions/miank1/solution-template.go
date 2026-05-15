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
	// TODO: Implement this function
	// 1. Clean the string (remove spaces, punctuation, and convert to lowercase)
	// 2. Check if the cleaned string is the same forwards and backwards
	
// lowercase
	s = strings.ToLower(s)

	// keep only letters and numbers
	var cleaned []rune

	for _, ch := range s {
		if unicode.IsLetter(ch) || unicode.IsNumber(ch) {
			cleaned = append(cleaned, ch)
		}
	}

	left := 0
	right := len(cleaned) - 1

	for left < right {
		if cleaned[left] != cleaned[right] {
			return false
		}

		left++
		right--
	}

	return true
}
