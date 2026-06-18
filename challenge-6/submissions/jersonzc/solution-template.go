// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	"strings"
	"unicode"
)

// getWords returns the words in a text.
func getWords(text string) []string {
	var answer []string

	text = strings.TrimSpace(text)

	var word []rune
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			word = append(word, r)
			continue
		}

		if unicode.IsSpace(r) || unicode.In(r, unicode.Hyphen) {
			if word != nil {
				answer = append(answer, string(word))
			}
			word = nil
		}
	}

	if word != nil {
		answer = append(answer, string(word))
	}

	return answer
}

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
	answer := make(map[string]int)

	text = strings.ToLower(text)
	words := getWords(text)

	for _, word := range words {
		if _, ok := answer[word]; ok {
			answer[word] += 1
		} else {
			answer[word] = 1
		}
	}

	return answer
}
