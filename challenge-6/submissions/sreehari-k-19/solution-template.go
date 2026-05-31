package challenge6

import (
	"strings"
	"unicode"
)

func CountWordFrequency(text string) map[string]int {
	result := make(map[string]int)
	cleanedText := strings.ToLower(text)
	cleanedText = strings.ReplaceAll(cleanedText, "'", "")

	words := strings.FieldsFunc(
		cleanedText,
		func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r)
		},
	)

	for _, word := range words {
		result[word]++
	}

	return result
}