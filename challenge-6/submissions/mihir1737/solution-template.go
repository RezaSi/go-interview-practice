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
// Output: map[string]int{"the": 2, "quick": 1, "brown": 1, "fox": 1, "jumps": 1,
// "over": 1, "lazy": 1, "dog": 1}

func normalizeString(s string) string {
	normalized := []rune("")

	for _, c := range s {
		if c >= '0' && c <= '9' {
			normalized = append(normalized, c)
		} else if c >= 'a' && c <= 'z' {
			normalized = append(normalized, c)
		} else if c >= 'A' && c <= 'Z' {
			normalized = append(normalized, c-'A'+'a')
		}
	}

	return string(normalized)
}
func CountWordFrequency(text string) map[string]int {
	splits := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '\''
	})

	ans := make(map[string]int)
	for _, s := range splits {
		normalized := normalizeString(s)
		ans[normalized]++
	}
	return ans
}
