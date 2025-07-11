// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	"bytes"
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
func CountWordFrequency(text string) map[string]int {	// Your implementation here
	// use whitespace and dash to parse out words
	f := func(c rune) bool {
		return unicode.IsSpace(c) || c == '-'
	}
	words := bytes.FieldsFunc([]byte(text), f)
	
	// initialise return map
	counts := make(map[string]int, len(words))
	
	// iterate through words filtering out invalid chars
	for _, word := range words {
		var b strings.Builder
		word = bytes.ToLower(word)
		for _, char := range word {
			if (char > 96 && char < 123) || (char > 47 && char < 58) {
				b.WriteByte(char)
			}
		}
		counts[b.String()]++
	}
	return counts
} 