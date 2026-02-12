// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	"strings"
)

// Package challenge6 contains the solution for Challenge 6.

// Add any necessary imports here

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
	if text == "" {
		return map[string]int{}
	}
	replacer := strings.NewReplacer(
		".", "", ",", "", "!", "", "?", "", ":", "", ";", "", "'", "", "-", " ",
		"\t", " ", "\n", " ",
	)
	text = replacer.Replace(text)
	text = strings.TrimSpace(text)
	text = strings.Join(strings.Fields(text), " ")

	wordArr := strings.Split(strings.ToLower(text), " ")

	wordMap := make(map[string]int)
	for _, word := range wordArr {
		wordMap[word]++
	}

	return wordMap
}
