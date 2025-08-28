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
func CountWordFrequency(text string) map[string]int {
	wordsSlice := strings.Fields(strings.ReplaceAll(text, "-", " "))
	re := regexp.MustCompile("[^A-Za-z0-9]*")
	wordFrequency := make(map[string]int, len(wordsSlice))
	for _, word := range wordsSlice {
		if word == "" {
			continue
		}
		lowWord := strings.ToLower(re.ReplaceAllString(word, ""))
		wordFrequency[lowWord]++
	}
	return wordFrequency
}

