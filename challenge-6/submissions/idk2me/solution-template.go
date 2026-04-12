// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	// Add any necessary imports here
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
	// Your implementation here
	m := make(map[string]int)
	words := strings.Fields(cleanText(text))
	for _, v := range words {
	    _, ok := m[v]
	    if ok {
	        m[v] += 1
	    } else {
	        m[v] = 1
	    }
	}
	return m
} 

func cleanText(s string) string {
    var b strings.Builder
    for _, v := range s {
        switch {
        case unicode.IsLetter(v) || unicode.IsNumber(v):
            b.WriteRune(unicode.ToLower(v))
        case unicode.IsSpace(v):
            b.WriteRune(' ')
        default:
            if v == '\'' {
                continue
            }
            b.WriteRune(' ')
        }
    }
    return b.String()
}





