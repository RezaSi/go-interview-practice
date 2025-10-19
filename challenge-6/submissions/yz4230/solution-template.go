// Package challenge6 contains the solution for Challenge 6.
package challenge6

import "strings"

func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '-'
}

func isLetter(c rune) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func CountWordFrequency(text string) map[string]int {
	// Your implementation here
	m := make(map[string]int)
	parts := make([]string, 0)

	buf := ""
	for _, c := range text {
		if len(buf) > 0 && isSpace(c) {
			parts = append(parts, buf)
			buf = ""
		} else if isLetter(c) {
			buf += string(c)
		}
	}
	if len(buf) > 0 {
		parts = append(parts, buf)
	}

	for _, part := range parts {
		m[strings.ToLower(part)]++
	}

	return m
}