package main

import (
	"fmt"
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
	var filtered string
	for _, r := range s {
	    if r >= rune('0') && r <= rune('9') {
	        filtered += string(r)
	    }
	    if r >= rune('a') && r <= rune('z') {
	        filtered += string(r)
	    }
	    if r >= rune('A') && r <= rune('Z') {
	        filtered += string(r - rune('A') + rune('a'))
	    }
	}
	for i := 0; i < len(filtered) / 2; i++ {
	    if filtered[i] != filtered[len(filtered) - i - 1] {
	        return false
	    }
	}
	return true
}
