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
	wordFreq := make(map[string]int)
	var sb strings.Builder

	for _, ch := range text {
	    if isAlphanumeric(ch) {
	        sb.WriteRune(ch)
	    } else {
	        if sb.Len() > 0 && (unicode.IsSpace(ch) || ch == '-') {
	            word := strings.ToLower(sb.String())
	            wordFreq[word]++
	            sb.Reset()
	        }
	    }
	}
	if sb.Len() > 0 {
	    word := strings.ToLower(sb.String())
        wordFreq[word]++
	}
	return wordFreq
} 

func isAlphanumeric(r rune) bool {
    return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}