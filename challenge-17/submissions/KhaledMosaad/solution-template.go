package main

import (
	"fmt"
	"regexp"
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
	re := regexp.MustCompile(`[\t\n\r-]`) // Matches tabs, newlines, and carriage returns
	cleanedString := re.ReplaceAllString(s, "")
	reg := regexp.MustCompile(`[^\p{L}\p{N} ]+`)
	text := reg.ReplaceAllString(cleanedString, "")
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "")
	for i, j := 0, len(text)-1; i <= j; {
		if text[i] != text[j] {
			return false
		}
		i++
		j--
	}

	return true
}
