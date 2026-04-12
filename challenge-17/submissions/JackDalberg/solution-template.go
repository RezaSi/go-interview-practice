package main

import (
	"fmt"
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
	s = strings.ToLower(s)
	
	var filtered strings.Builder
	for _, r := range s {
	    if '0' <= r && r <= '9' || 'a' <= r && r <= 'z' {
	        filtered.WriteRune(r)
	    }
	}
	filteredStr := filtered.String()
	length := len(filteredStr)
	for i := 0; i < length/2; i++ {
	    if filteredStr[i] != filteredStr[length-1-i]{
	        return false
	    }
	}
	return true
}

