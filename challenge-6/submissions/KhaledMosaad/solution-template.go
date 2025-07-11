// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	"regexp"
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

func getPureWords(s string) []string {
    re := regexp.MustCompile(`[\t\n\r-]`) // Matches tabs, newlines, and carriage returns
	cleanedString := re.ReplaceAllString(s, " ")
	reg := regexp.MustCompile(`[^\p{L}\p{N} ]+`)
	text := reg.ReplaceAllString(cleanedString, "")
	
	lowercaseText := strings.ToLower(text)
	words := strings.Split(lowercaseText, " ")
	return words
}

func CountWordFrequency(text string) map[string]int {
    words := getPureWords(text)
	ret := make(map[string]int)
	for _, word := range words {
	    if word != "" {
	        ret[word]++
	    }	    
	}
	return ret
} 