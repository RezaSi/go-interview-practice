// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	// Add any necessary imports here
	"strings"
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
    result := make(map[string]int)
    if len(text) == 0  {
        return result
    }
	
	re := regexp.MustCompile(`[^a-zA-Z0-9']+`)
    text = strings.ReplaceAll(text, "'", "")
    // Replace all non-alphanumeric characters with a space
    cleaned := re.ReplaceAllString(strings.ToLower(text), " ")  // "Hello world 123"
    
    // Split using regex
    words := re.Split(strings.TrimSpace(cleaned), -1)
    
    for _, word := range words {
        _, exists := result[word]
        if exists {
            result[word] += 1
        } else {
            result[word] = 1
        }
    }
	return result
} 