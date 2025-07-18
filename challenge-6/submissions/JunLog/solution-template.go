// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	// Add any necessary imports here
	"strings"
	"unicode"
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
    text = strings.Replace(text,`'`,"",-1)
	f := func(c rune)bool{
	    return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	resMap := map[string]int{}
	split := strings.FieldsFunc(text,f)
	for _,_w := range split{
	    if string(_w) == ""{
	        continue
	    }
	    w := strings.ToLower(_w)
	    if _,ok := resMap[w];!ok{
	        resMap[w] = 0
	    }
	    resMap[w] += 1
	}
	return resMap
} 