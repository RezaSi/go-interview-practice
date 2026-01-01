package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func IsPalindrome(s string) bool {
	var builder strings.Builder
	for _, c := range s {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			builder.WriteRune(unicode.ToLower(c))
		}
	}
	normalized := builder.String()
	left, right := 0, len(normalized)-1
	for left < right {
		if normalized[left] != normalized[right] {
			return false
		}
		left++
		right--
	}
	return true
}
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	result := IsPalindrome(input)
	if result {
		fmt.Println("The string is a palindrome.")
	} else {
		fmt.Println("The string is not a palindrome.")
	}
}