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
    runes := []rune(s)
    left, right := 0, len(runes)-1

    for left < right {
        // Пропускаем слева всё, что не буква и не цифра
        if !unicode.IsLetter(runes[left]) && !unicode.IsDigit(runes[left]) {
            left++
            continue
        }
        // Пропускаем справа всё, что не буква и не цифра
        if !unicode.IsLetter(runes[right]) && !unicode.IsDigit(runes[right]) {
            right--
            continue
        }

        // Сравниваем, приведя к нижнему регистру
        if unicode.ToLower(runes[left]) != unicode.ToLower(runes[right]) {
            return false
        }

        left++
        right--
    }
    return true
}
