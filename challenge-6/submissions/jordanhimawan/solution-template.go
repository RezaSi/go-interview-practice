// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	// Add any necessary imports here
	"strings"
	"fmt"
	"regexp"
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

	sentence := strings.ToLower(text)
	sentence = strings.ReplaceAll(sentence, "'", "")  
	reg := regexp.MustCompile(`\b[a-zA-Z0-9]+\b`)

	words := reg.FindAllString(sentence, -1)
	fmt.Println("words: ", words)
	
	wordFreq := make(map[string]int)
	for _, word := range words {
	    if _, exists := wordFreq[word]; exists {
	        wordFreq[word] += 1
	        continue
	    }
	    
	    wordFreq[word] = 1
	}
	
	return wordFreq
} 