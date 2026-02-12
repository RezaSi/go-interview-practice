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
	lowered := strings.ToLower(s)
	
	var cleansed strings.Builder
	for _, ch := range lowered {
	    if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') {
	        cleansed.WriteRune(ch)
	    }
	}
	
	result := cleansed.String()
	
	// 2. Check if the cleaned string is the same forwards and backwards
	for i := 0; i < len(result)/2; i++ {
	    if result[i] != result[len(result)-1-i] {
	        return false
	    }
	}
	return true
}
