package main

import (
	"fmt"
	"strings"
	"unicode"
)

func main() {
	var input string
	fmt.Scanln(&input)
	result := IsPalindrome(input)
	if result {
		fmt.Println("The string is a palindrome.")
	} else {
		fmt.Println("The string is not a palindrome.")
	}
}

func IsPalindrome(s string) bool {
	var builder strings.Builder
	for _, c := range s {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			builder.WriteRune(unicode.ToLower(c))

		}
	}
	sClean := builder.String()
	n := len(sClean)
	for i := 0; i < n/2; i++ {
		if sClean[i] != sClean[n-i-1] {
			return false
		}
	}
	return true
}