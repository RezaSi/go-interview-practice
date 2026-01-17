// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
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
	wordsFreq := make(map[string]int, int(len(text) / 10))
	
	replacedText := strings.ReplaceAll(text, "'", "")
	
	var i, j int
	for idx := 0; idx < len(replacedText); idx++ {
	    if isDigit(replacedText[idx]) || isLatin(replacedText[idx]) {
	        j++
	    } else {
	        if j - i > 0 {
	            word := strings.ToLower(replacedText[i:j])
	            wordsFreq[word]++
	        }
	        i, j = idx + 1, idx + 1
	    }
	}
	
	if j - i > 0 {
	    word := strings.ToLower(replacedText[i:j])
        wordsFreq[word]++
	}
	
	return wordsFreq
} 

func isDigit(c byte) bool {
    return c >= '0' && c <= '9'
}

func isLatin(c byte) bool {
    return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == 27
}
