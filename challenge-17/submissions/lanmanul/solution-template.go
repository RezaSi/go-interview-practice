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
    // Нормализация строки: нижний регистр + удаление всех не-буквенных и не-цифровых символов
    var cleaned strings.Builder
    for _, r := range strings.ToLower(s) {
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            cleaned.WriteRune(r)
        }
    }
    cleanedStr := cleaned.String()
    
    runes := []rune(cleanedStr)
    for i := 0; i < len(runes)/2; i++ {
        if runes[i] != runes[len(runes)-i-1] {
            return false
        }
    }
    return true
}