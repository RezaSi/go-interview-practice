package main

import (
	"fmt"
	"strings"
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
// A palindrome reads the same backward as forward, ignoring case, 
// spaces, and punctuation.
func IsPalindrome(s string) bool {
    normalizeString := normalizeString(s)
    
    left := 0
    right := len(normalizeString)-1
    
    for left < right {
        if normalizeString[left] != normalizeString[right] {
            return false
        }
        left++
        right--
    }
    
	return true
}

func normalizeString(str string) string {
   res := ""

    for _, s := range str {
        if ('a' <= s && s <= 'z') || 
           ('A' <= s && s <= 'Z') || 
           ('0' <= s && s <= '9') {
            res += string(s)
        }
    }
    
    
    
    return strings.ToLower(res)
}


