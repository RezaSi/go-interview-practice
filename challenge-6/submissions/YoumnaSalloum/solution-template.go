// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
"strings"
"unicode"
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
	// Step 1: Clean the text by turning any non-letter/non-digit into a space, 
	// and convert everything to lowercase on the fly.
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		if r == '\'' || r == '’' {
			return -1 // Returning a negative value drops the character in strings.Map
		}
		return ' ' // Replace punctuation, hyphens, and symbols with spaces
	}, text)

	// Step 2: Split the cleaned text into words by spaces
	sp := strings.Fields(cleaned)
	
	// Step 3: Count word frequencies
	words := make(map[string]int)
	for _, w := range sp {
		words[w]++
	}

	return words
}