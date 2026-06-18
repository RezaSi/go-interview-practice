// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	"strings"
	"unicode"
)

func CountWordFrequency(text string) map[string]int {
	var cleaned strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cleaned.WriteRune(unicode.ToLower(r))
		} else if r == '\'' { 
			continue
		} else {
			cleaned.WriteRune(' ')
		}
	}

	words := strings.Fields(cleaned.String())
	counter := make(map[string]int)

	for _, word := range words {
		counter[word]++
	}

	return counter
}

