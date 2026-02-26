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

	r_list := []rune(s)

	i, j := 0, len(r_list)-1

	for i < j {
		i_v := r_list[i]
		if !(unicode.IsLetter(i_v) || unicode.IsDigit(i_v)) {
			i += 1
			continue
		}

		j_v := r_list[j]
		if !(unicode.IsLetter(j_v) || unicode.IsDigit(j_v)) {
			j -= 1
			continue
		}

		if unicode.ToLower(i_v) != unicode.ToLower(j_v) {
			return false
		}

		i += 1
		j -= 1

	}
	return true
}