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
// Words like "let's" convert to "lets"
// All words are converted to lowercase before counting.
// All punctuation, spaces, and other non-alphanumeric characters are ignored.
//
// Example:
// Input: "The quick brown fox jumps over the lazy dog."
// Output: map[string]int{"the": 2, "quick": 1, "brown": 1, "fox": 1, "jumps": 1, "over": 1, "lazy": 1, "dog": 1}
func CountWordFrequency(text string) map[string]int {

	// Define new map res to return
	res := make(map[string]int)

	// Normalize input string
	lowerString := strings.ToLower(text)
	lowerString = strings.ReplaceAll(lowerString, "'s", "s")

	// Replace every non-alphanumeric character with a space.
	repString := strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return ' '
	}, lowerString)

	// Loop of the words in repString and add counter to res
	for _, s := range strings.Fields(repString) {
		res[s]++
	}
	return res
}

