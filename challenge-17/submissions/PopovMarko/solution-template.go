package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	// Get input from the user
	fmt.Print("Enter a string to check if it's a palindrome: ")

	// New reader with Stdin returns buffered reader
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	// input string normalization
	input = strings.TrimRight(input, "\r\n")

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
	if s == "" {
		return true
	}

	// Normalization of the string to lower case
	s = strings.ToLower(s)

	//delete all not alphanumeric runes
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, s)

	// Get a slice of runes to keep all UTF-8 symbol
	runes := []rune(s)
	l := 0
	r := len(runes) - 1

	// Loop with two indexes left and right that steps towards each other
	// and compare each runes
	for l <= r {
		if runes[l] != runes[r] {
			return false
		}
		l++
		r--
	}
	return true
}
