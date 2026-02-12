package challenge6

import (
	"strings"
	"unicode"
)

// CountWordFrequency takes a string containing multiple words and returns
// a map where each key is a word and the value is the number of times that
// word appears in the string. The comparison is case-insensitive.
func CountWordFrequency(text string) map[string]int {
	frequency := make(map[string]int)

	if text == "" {
		return frequency
	}

	var currentWord strings.Builder

	for _, char := range text {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			currentWord.WriteRune(unicode.ToLower(char))
		} else if char == '\'' {
			// ignore apostrophes (do nothing)
			continue
		} else {
			// word boundary
			if currentWord.Len() > 0 {
				word := currentWord.String()
				frequency[word]++
				currentWord.Reset()
			}
		}
	}

	// handle last word if exists
	if currentWord.Len() > 0 {
		word := currentWord.String()
		frequency[word]++
	}

	return frequency
}
