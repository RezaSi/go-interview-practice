// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
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
	fmap := map[string]int{}
	text = strings.ToLower(text)
	// split the string into case-insensitive words
	// while splitting only consider alphabets and numbers to included as part of word, rest all characters as delimiters
	words := []string{}
	var temp string
	
	for i:=0 ; i<len(text); i++ {
	    if (text[i] >= 'a' && text[i]<='z') || (text[i] >= '0' && text[i]<='9') {
	        temp += string(text[i])
	    } else if (text[i]=='\\') {
	        i++
	    } else if (text[i]=='\'') {
	        continue
	    } else if len(temp) > 0 {
	        words = append(words, temp)
	        temp = ""
	    }
	}
	if len(temp) > 0 {
	    words = append(words, temp)
	}
	
	for _, w := range words {
	    _, ok := fmap[w]
	    if ok {
	        fmap[w] += 1
	    } else {
	        fmap[w] = 1
	    }
	}
	return fmap
} 