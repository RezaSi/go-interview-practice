// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
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
	count := make(map[string]int)
	runes := []rune(text)

	// average word length in English is around 5 characters,
	// so we can use that as an initial capacity hint
	buf := make([]rune, 0, 5)
	for _, r := range runes {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			buf = append(buf, unicode.ToLower(r))
		} else if r == '\'' {
			continue
		} else {
			if len(buf) > 0 {
				count[string(buf)]++
				buf = buf[:0]
			}
		}
	}

	if len(buf) > 0 {
		count[string(buf)]++
	}

	return count
}
