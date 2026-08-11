// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	// Add any necessary imports here
	"regexp"
	"strings"
)

// CountWordFrequency takes a string containing multiple words and returns
// a map where each key is a word and the value is the number of times that
// word appears in the string. The comparison is case-insensitive.
//
// Words are defined as sequences of letters and digits.
// All words are converted to lowercase before counting.
// All punctuation, spaces, and other non-alphanumeric characters are ignored.
//
// For example:
// Input: "The quick brown fox jumps over the lazy dog."
// Output: map[string]int{"the": 2, "quick": 1, "brown": 1, "fox": 1, "jumps": 1, "over": 1, "lazy": 1, "dog": 1}
func CountWordFrequency(text string) map[string]int {
	// Your implementation here
	response := make(map[string]int)
	// step 1 : split with regex to exclude special character
	
	text = strings.ReplaceAll(text, "-", " ")
	x := regexp.MustCompile(`[^a-zA-Z0-9 \s]+`)
	arr := strings.ToLower(x.ReplaceAllString(text, ""))
	arr = strings.ReplaceAll(arr, "\n", " ")
	arr = strings.ReplaceAll(arr, "\t", " ")
	arr = strings.ReplaceAll(arr, "\r", " ")

	// fmt.Println(arr + " return from regexReplaceSpecialChar")
	split := strings.Split(arr, " ")
	for _, value := range split {
		// step 2 : check if value is exist in key response
		if value == "" {
			continue
		} else {
			valuee, exists := response[value]
			if exists {
				response[value] = valuee + 1
			} else {
				response[value] = 1
			}
		}
	}
	return response
} 