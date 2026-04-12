package challenge6

import (
	"strings"
	"unicode"
)

func CountWordFrequency(text string) map[string]int {
	word := make(map[string]int)
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, "'", "")

	var builder strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		} else {
			builder.WriteRune(' ')
		}
	}

	words := strings.Fields(builder.String())
	for _, w := range words {
		word[w]++
	}
	return word
}
