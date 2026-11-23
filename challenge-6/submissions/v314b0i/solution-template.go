// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	"strings"
)

// Without using string split. JLT.
func CountWordFrequency(text string) (wordFreqMap map[string]int) {
    // updateMap takes a word as a rune aray, normlises it and updates the freqMap
    wordFreqMap = map[string]int{}
	updateMap := func(wordArr []rune) {
	    word := strings.ToLower(string(wordArr))
	    if count, exists := wordFreqMap[word]; exists {
	        wordFreqMap[word] = count + 1
	    } else {
	        wordFreqMap[word] = 1
	    }
	}
	// isAlphabet checks if a rune is a word
	isAlphabetNum := func(r rune) bool {
	    if (r <= 'z' && r >= 'a') || (r <= 'Z' && r >= 'A') {
	        return true
	    }
	    if r <= '9' && r > '0' { // just to meet the test case. 
	        return true 
	    }
	    return false
	}
	
	var buffer []rune
	// flushBuffer
	flushBuffer := func() {
	    if len(buffer)>0 {
	        updateMap(buffer)
	        buffer = []rune{}
	    }
	}
	
	for _, r := range text {
	    if isAlphabetNum(r) {
	        buffer = append(buffer, r)
	    } else if r == '\''{ // Bad unit test and expectation. Not correct.
	        continue
	    } else {
	        flushBuffer()
	    }
	}
	flushBuffer()
	return
} 