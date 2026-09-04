// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	"strings"
	"unicode"
)

func CountWordFrequency(text string) map[string]int {
	hashMap := make(map[string]int)
	
	words := strings.FieldsFunc(text, func(c rune) bool {
	    return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '\''
	})
	
	for _, word := range words {
	    parsedWord := removeNonAlphabeticChar(strings.ToLower(word))
	    if parsedWord == "" {
	        continue
	    }
	    
	    hashMap[parsedWord]++
	}
	
	return hashMap
}

func removeNonAlphabeticChar(text string) string {
    var newText strings.Builder
    
    for _, c := range text {
        if isAlphaNumeric(c) {
            newText.WriteRune(c)
        }
    }
    
    return newText.String()
}

func isAlphaNumeric(c rune) bool {
    if (c >= 'a' && c <= 'z')  || (c >= '0' && c <= '9') {
        return true
    }
    return false
}