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
	s = strings.ToLower(s)

	re := regexp.MustCompile(`[^a-z0-9]`)

	cleared := re.ReplaceAllString(s, "")

	left := 0
	right := len(cleared) - 1

	for left < right {
		if cleared[left] != cleared[right] {
			return false
		}

		left++
		right--
	}

	return true
}