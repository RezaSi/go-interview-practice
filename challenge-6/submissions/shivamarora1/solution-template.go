// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
    "strings"
    "regexp"
	// Add any necessary imports here
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
	// Your implementation here
	text = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(text, "\t", " "), "\n", " "),"-"," ")
	reg := regexp.MustCompile("[^a-zA-Z0-9 ]+")
	cleanString := reg.ReplaceAllString(text, "")
    
    words := strings.Split(cleanString," ")
	result := make (map[string]int)
	for _,word := range words{
	    if word == "" {
			continue
		}
	    lWord:= strings.ToLower(string(word))
	    if _,ok := result[lWord];ok{
	        result[lWord]+=1
	    }else{
	        result[lWord] = 1
	    }
	}
	return result
} 