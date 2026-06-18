// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	"regexp"
	"strings"
)

var nonAlphaNumericRegex = regexp.MustCompile(`[\W_]+`)
var splitRegex = regexp.MustCompile(`[\s-]+`)

func CountWordFrequency(text string) map[string]int {
	frequency := map[string]int{}

	for _, word := range splitRegex.Split(text, -1) {
		lowerWord := strings.ToLower(word)
		cleanWord := nonAlphaNumericRegex.ReplaceAllLiteralString(lowerWord, "")
		if cleanWord == "" {
			continue
		}

		frequency[cleanWord]++
	}

	return frequency
}
