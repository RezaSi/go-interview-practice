// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
    "regexp"
    "strings"
)

var (
    alphaNumeric = regexp.MustCompile("[^A-Za-z0-9\\s]")
    splitRegex = regexp.MustCompile("\\s+|-")
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
    counts := make(map[string]int)
    
    if len(text) == 0 {
        return counts
    }
    
    text = string(splitRegex.ReplaceAll([]byte(text), []byte(" ")))
    text = string(alphaNumeric.ReplaceAll([]byte(text), []byte("")))
    
    words := strings.Fields(strings.ToLower(text))
    
    for _, w := range words {
        if _, ok := counts[w]; !ok {
            counts[w] = 0
        }
        
        counts[w]++
    }
    
	return counts
} 