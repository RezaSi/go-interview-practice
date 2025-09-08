// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func clearString(str string) string {
    s := strings.Replace(str, "'", "", -1)
    s = strings.ToLower(nonAlphanumericRegex.ReplaceAllString(s, " "))
    return s
}

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

    
	cleared := clearString(text)
	//fmt.Printf(cleared)
	
	s := strings.Split(cleared, " ")
	
	for _, word := range s {
	    if word != "" {
    	    cnt, ok := result[word]    
    	    
    	    if ok {
    	        result[word] = cnt+1
    	    } else {
    	        result[word] = 1
    	    }
	    }
	}
	
	return result
} 