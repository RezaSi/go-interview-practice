// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	"strings"
	"unicode"
)


func CountWordFrequency(text string) map[string]int {

	lower := strings.Map(func(c rune) rune {
		if unicode.IsPunct(c) && c != '-' {
			return -1
		} else if c == '-' {
			return ' '
		}
		return c
	}, strings.ToLower(text))

	splitter := strings.FieldsFunc(lower, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
	
	wCount := make(map[string]int)
	for _, w := range splitter {

		if wCount[w] == 0 {
			wCount[w] = 1
		} else {
			wCount[w]++
		}
	}
	return wCount
} 