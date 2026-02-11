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
	var clearStr string
	for _, v := range s {
		if rune(v) >= 65 && rune(v) <= 90 {
			clearStr += string(v + 32)
		}
		if (rune(v) >= 97 && rune(v) <= 122) || (rune(v) >= 48 && rune(v) <= 57) {
			clearStr += string(v)
		}
	}

	fmt.Println(clearStr[:(len(clearStr)+1)/2])
	fmt.Println(clearStr[len(clearStr)/2:])

	var reverse string

	for _, v := range clearStr[len(clearStr)/2:] {
		reverse = string(v) + reverse
	}

	return reverse == clearStr[:(len(clearStr)+1)/2]
}
