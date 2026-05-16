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
    cleaned := []rune{}
    
    for _, r := range s{
        if unicode.IsLetter(r) || unicode.IsDigit(r){
            lower := unicode.ToLower(r)
            cleaned = append(cleaned, lower)
        }
    }
    
    i, j := 0, len(cleaned) - 1
    
    for i < j {
        if (cleaned[i] != cleaned[j]){
            return false
        }
        
        i = i + 1
        j = j - 1
    }
	return true
}
