package main

import (
	"fmt"
	"unicode"
)

func main() {

	// Call the IsPalindrome function and print the result
	fmt.Println(IsPalindrome("!@#$%^&*()😊1221😊dFSFDvfdv"))
	fmt.Println(IsPalindrome("Madam, I'm Adam"))

}

// IsPalindrome checks if a string is a palindrome.
// A palindrome reads the same backward as forward, ignoring case, spaces, and punctuation.
func IsPalindrome(s string) bool {
	// TODO: Implement this function
	// 1. Clean the string (remove spaces, punctuation, and convert to lowercase)
	// 2. Check if the cleaned string is the same forwards and backwards
	ps := prepareText(s)

	l := 0
	r := len(ps) - 1

	for l <= r {
		if ps[l] != ps[r] {
			return false
		}

		l += 1
		r -= 1
	}

	return true
}

func prepareText(input string) string {
	buf := make([]rune, 0, len(input))
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			buf = append(buf, unicode.ToLower(r))
		}
	}
	return string(buf)
}
