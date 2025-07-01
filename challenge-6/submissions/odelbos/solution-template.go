// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
    "regexp"
    "strings"
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
    result := make(map[string]int)
	re := regexp.MustCompile(`[\w']+`)
	for _, word := range(re.FindAllString(text, -1)) {
	    lower := strings.ReplaceAll(strings.ToLower(word), "'", "")
	    count, ok := result[lower]
	    if ok {
	        result[lower] = count + 1
	    } else {
	        result[lower] = 1
	    }
	}
	return result
}