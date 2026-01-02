// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
    "strings"
    "regexp"
)

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

	m := make(map[string]int, 0)
	
	if text == "" {
		return m
	}
	
	text = strings.ToLower(text)
	cutOff := []string{"'", "", "\t", " ", "\n", " "}
	text = strings.NewReplacer(cutOff...).Replace(text)

	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	cleaned := re.ReplaceAllString(text, " ")
	cleaned = strings.Trim(cleaned, " ")



	arr := re.Split(cleaned, -1)

	for _, value := range arr {

		_, exist := m[value]
		if exist {
			m[value]++
			continue
		}
		m[value] = 1
	}

	return m
}