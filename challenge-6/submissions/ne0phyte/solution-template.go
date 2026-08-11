// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	"strings"
	"unicode"
)

// Add any necessary imports here

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
	words := make(map[string]int)
	for _, token := range tokenize(text) {
		if len(token) > 0 {
			words[token]++
		}
	}
	return words
}

func tokenize(text string) []string {
	words := []string{}
	var token strings.Builder
	for _, r := range []rune(text) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' {
			token.WriteRune(unicode.ToLower(r))
		} else if r != '\'' {
			words = append(words, token.String())
			token.Reset()
		}
	}
	return append(words, token.String())
}
