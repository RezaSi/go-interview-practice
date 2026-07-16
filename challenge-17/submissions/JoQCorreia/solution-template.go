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
    
    //Removes all characters that are not letters or numbers and makes string lowercase
	pal := strings.Map(func(c rune) rune {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			return -1
		}
		return c
	}, strings.ToLower(s))

    //If string has no characters, it returns true
    if len(pal) == 0 {
		return true
	}

    // Convert to runes to safely handle multi-byte Unicode characters
	palRunes := []rune(pal)
	indx := len(palRunes) - 1

	for i := 0; i <= (indx - i); i++ {
    //Loop compares the letters on opposite sides of the string for matches
		if palRunes[i] != palRunes[indx-i] {
			return false
		}

		continue

	}

	return true
}