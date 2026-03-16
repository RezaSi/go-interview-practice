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
// IsPalindrome checks if a string is a palindrome.
// A palindrome reads the same backward as forward, ignoring case, spaces, and punctuation.
func IsPalindrome(s string) bool {
    // Clean the string: convert to lowercase and remove non-alphanumeric characters
    cleaned := make([]rune, 0, len(s))
    for _, r := range s {
        if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
            cleaned = append(cleaned, unicode.ToLower(r))
        }
    }
    
    // Check if the cleaned string is a palindrome
    length := len(cleaned)
    for i := 0; i < length/2; i++ {
        if cleaned[i] != cleaned[length-1-i] {
            return false
        }
    }
    
    return true
}