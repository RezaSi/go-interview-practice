package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	// Get input from the user
	input := "A man, a plan, a canal: Panama"
	// Call the IsPalindrome function and print the result
	result := IsPalindrome(input)
	fmt.Println(result)
}

// IsPalindrome checks if a string is a palindrome.
// A palindrome reads the same backward as forward, ignoring case, spaces, and punctuation.
func IsPalindrome(s string) bool {
	re := regexp.MustCompile(`\W`)
	cleanedText := strings.ToLower(re.ReplaceAllString(s, ""))
	fmt.Println(cleanedText)
	for i, j := 0, len(cleanedText)-1; i < j; i, j = i+1, j-1 {
		if cleanedText[i] != cleanedText[j] {
			return false
		}
	}
	return true
}
