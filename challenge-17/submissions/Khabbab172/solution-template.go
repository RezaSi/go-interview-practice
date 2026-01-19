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
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	cleaned := re.ReplaceAllString(s,"")
	cleaned = strings.ToLower(cleaned)
	start := 0
	end := len(cleaned) - 1
	for start <= end {
	    if cleaned[start] != cleaned[end] {
	        return false
	    }
	    start++
	    end--
	}
	return true
}
