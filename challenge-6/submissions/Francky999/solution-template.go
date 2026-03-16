package challenge6

import (
	"strings"
	"regexp"
)
var (
	nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9\s]|\t|\n`)
	multiSpaceRegex      = regexp.MustCompile(`\s{2,}`)
)

func CountWordFrequency(text string) map[string]int {
	text = strings.ReplaceAll(text, "'", "")
	replaced := nonAlphanumericRegex.ReplaceAll([]byte(strings.ToLower(strings.TrimSpace(text))), []byte(" "))
	replaced = multiSpaceRegex.ReplaceAll(replaced, []byte(" "))
	words, wordsCounter := strings.Split(string(replaced), " "), map[string]int{}
	for _, word := range words {
		wordsCounter[word] += 1
	}
	delete(wordsCounter, "")
	return wordsCounter
}
