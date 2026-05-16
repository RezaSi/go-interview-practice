package main

import (
	"fmt"
	"strings"
	"regexp"
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
	// 2. Check if the cleaned string is the same forwards and backwards'
	if s==""{
	    return true
	}
temp:=strings.ToLower(s)
	reg := regexp.MustCompile("[^a-z0-9]")
	temp= reg.ReplaceAllString(temp, "")
    	left, right := 0, len(temp)-1

	for left < right {
		if temp[left] != temp[right] {
			return false
		}
		left++
		right--
	}
	return true
}
