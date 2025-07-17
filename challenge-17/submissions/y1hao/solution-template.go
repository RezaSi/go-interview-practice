package main

import (
	"fmt"
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
	runes := []rune(s)
	i, j := 0, len(runes)-1
	for i < j {
		for i < len(runes) && !isValid(runes[i]) {
			i++
		}
		for j >= 0 && !isValid(runes[j]) {
			j--
		}
		if i >= j || i == len(runes) || j < 0 {
			return true
		}
		if unicode.ToLower(runes[i]) != unicode.ToLower(runes[j]) {
			return false
		}
		i++
		j--
	}
	return true
}

func isValid(r rune) bool {
	return !unicode.IsSpace(r) && !unicode.IsPunct(r) && !unicode.IsSymbol(r)
}
