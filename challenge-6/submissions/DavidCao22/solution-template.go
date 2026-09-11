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
// Hyphenated words should be split into separate words.
//
// For example:
// Input: "The quick brown fox jumps over the lazy dog."
// Output: map[string]int{"the": 2, "quick": 1, "brown": 1, "fox": 1, "jumps": 1, "over": 1, "lazy": 1, "dog": 1}
func CountWordFrequency(text string) map[string]int {
	whiteSpaceRegex := regexp.MustCompile(`[\s]+`)
	nonLettersAndNumsRegex := regexp.MustCompile(`[^A-Za-z\d -]+`)
	hyphenRegex := regexp.MustCompile(`-`)
	
	whiteSpaceNowSpaces := whiteSpaceRegex.ReplaceAllString(text, " ")
	allLowerCase := strings.ToLower(whiteSpaceNowSpaces)

	allLettersAndNums := nonLettersAndNumsRegex.ReplaceAllString(allLowerCase, "")
	separatedHyphenatedWords := hyphenRegex.ReplaceAllString(allLettersAndNums, " ")

	words := strings.Split(separatedHyphenatedWords, " ")

	count := make(map[string]int)
	for _, w := range words {
	    if w != "" {
    	    if _, ok := count[w]; ok {
    	        count[w] = count[w] + 1
    	    } else {
    	        count[w] = 1
    	    }
	    }
	}
	
	return count
} 