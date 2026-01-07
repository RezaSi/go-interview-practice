package challenge6

import (
	"fmt"
	"strings"
	"unicode"
)

func CountWordFrequency(text string) map[string]int {
	wordsCount := make(map[string]int, 1000) // Pre-allocate for 1000 words
	var builder strings.Builder
	for _, c := range text {
		c = unicode.ToLower(c)
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			builder.WriteRune(c)
			// builder.WriteString(string(c))
		} else if c != '\'' {
			if builder.Len() > 0 {
				wordsCount[builder.String()] += 1
				builder.Reset()
			}
		}
	}
	if builder.Len() > 0 {
		wordsCount[builder.String()] += 1
	}
	return wordsCount
}

func main() {
	fmt.Println(CountWordFrequency("Let's go baby my baby"))
}
