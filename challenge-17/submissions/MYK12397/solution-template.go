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
	// TODO: Implement this function
	// 1. Clean the string (remove spaces, punctuation, and convert to lowercase)
	// 2. Check if the cleaned string is the same forwards and backwards
    cleanString := strings.ReplaceAll(s, "'", "")
    cleanString = strings.ReplaceAll(cleanString, "'", "")
    res := ""
    str  := strings.Map(func(r rune) rune {
        if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r){
            return r
        }
        return ' '
    }, cleanString)
    str = strings.ToLower(str)
    for _,w := range str {
        if string(w) != " " {
        res += string(w)
        }
    }
    fmt.Println("str: ",res)
    for i := 0; i < len(res); i++ {
        if res[i] != res[len(res)-1-i] {
            return false
        }
    }
	return true
}
