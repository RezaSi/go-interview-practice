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
    cleanText := strings.ReplaceAll(text, "'", "")
    cleanText = strings.ReplaceAll(cleanText, "'", "") 
	mp := map[string]int{}
 cleanSentence := strings.Map(func(r rune) rune {
    if unicode.IsLetter(r) || unicode.IsSpace(r) || unicode.IsNumber(r){
        return r
    }
    return ' '
   }, cleanText)

	words := strings.Fields(cleanSentence)
	for _, word := range words {
	    if word != "" {
	    mp[strings.ToLower(string(word))]++
	    
	    }
	}
	return mp
} 