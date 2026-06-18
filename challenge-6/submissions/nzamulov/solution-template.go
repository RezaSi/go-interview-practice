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
func validSymbol(r rune) bool {
    return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func spaceSymbol(r rune) bool {
    return r == ' ' || r == '\t' || r == '\n' || r == '-'
}

func CountWordFrequency(text string) map[string]int {
	var current strings.Builder
	m := make(map[string]int)
	for _, r := range text {
	    if !spaceSymbol(r) {
	        if validSymbol(r) {
                current.WriteRune(unicode.ToLower(r))   
	        }
	    } else if current.Len() > 0 {
	        m[current.String()]++
	        current.Reset()
	    }
	}
	if current.Len() > 0 {
	    m[current.String()]++
	}
	return m
} 