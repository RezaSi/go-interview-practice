// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	// Add any necessary imports here
	"unicode"
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
	// handle Empty string
	if len(text) == 0 {
	    return make(map[string]int)
	}
	// handle text: lowercase and add period
	s := strings.ToLower(text)
	l := rune(s[len(s)-1])
	if unicode.IsLetter(l) || unicode.IsDigit(l) {
		tmp := []rune(s)
		tmp = append(tmp, '.')
		s = string(tmp)
	}
	word := []rune{}
	collector := make(map[string]int)
	for _, v := range s {
	    if unicode.IsLetter(v) || unicode.IsDigit(v) {
	        word = append(word, v)
	    } else if v == 39 {
	        continue  
	    } else {
	        if len(word) > 0 {
	            collector[string(word)] += 1
	        }
	        word = []rune{}
	    }
	} 
	return collector
} 