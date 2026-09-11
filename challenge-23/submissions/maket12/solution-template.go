package main

import (
	"fmt"
)

func main() {
	// Sample texts and patterns
	testCases := []struct {
		text    string
		pattern string
	}{
		{"ABABDABACDABABCABAB", "ABABCABAB"},
		{"AABAACAADAABAABA", "AABA"},
		{"GEEKSFORGEEKS", "GEEK"},
		{"AAAAAA", "AA"},
	}

	// Test each pattern matching algorithm
	for i, tc := range testCases {
		fmt.Printf("Test Case %d:\n", i+1)
		fmt.Printf("Text: %s\n", tc.text)
		fmt.Printf("Pattern: %s\n", tc.pattern)

		// Test naive pattern matching
		naiveResults := NaivePatternMatch(tc.text, tc.pattern)
		fmt.Printf("Naive Pattern Match: %v\n", naiveResults)

		// Test KMP algorithm
		kmpResults := KMPSearch(tc.text, tc.pattern)
		fmt.Printf("KMP Search: %v\n", kmpResults)

		// Test Rabin-Karp algorithm
		rkResults := RabinKarpSearch(tc.text, tc.pattern)
		fmt.Printf("Rabin-Karp Search: %v\n", rkResults)

		fmt.Println("------------------------------")
	}
}

// NaivePatternMatch performs a brute force search for pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func NaivePatternMatch(text, pattern string) []int {
    if len(text) == 0 || len(pattern) == 0 || len(pattern) > len(text) {
        return make([]int, 0)
    }
    
    result := make([]int, 0)
    
    for i := 0; i <= len(text) - len(pattern); i++ {
        var j int
        
        for j < len(pattern) && text[i+j] == pattern[j] {
            j++
        }
        
        if j == len(pattern) {
            result = append(result, i)
        }
    }
    
    return result
}

func countLPS(pattern string) []int {
    lps := make([]int, len(pattern))
    
    var i, suffLen = 1, 0
    
    for i < len(pattern) {
        if pattern[i] == pattern[suffLen] {
            suffLen++
            lps[i] = suffLen
            i++
        } else {
            if suffLen != 0 {
                suffLen = lps[suffLen-1]
            } else {
                lps[i] = 0
                i++
            }
        }
    }
    
    return lps
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	result := make([]int, 0)
	
	if len(text) == 0 || len(pattern) == 0 || len(text) < len(pattern) {
	    return result
	}
	
	lps := countLPS(pattern)
	
	var i, j int = 0, 0
	
	for i < len(text) {
	    if pattern[j] == text[i] {
	        i++
	        j++
	    }
	    
	    if j == len(pattern) {
	        result = append(result, i-j)
	        j = lps[j-1]
	    } else if i < len(text) && pattern[j] != text[i] {
	        if j != 0 {
	            j = lps[j-1]
	        } else {
	            i++
	        }
	    }
	}
	
	return result
}

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	result := make([]int, 0)
	
	if len(text) == 0 || len(pattern) == 0 || len(text) < len(pattern) {
	    return result
	}
	
	// Hash counting
	const prime, base = 101, 256
	
	var patternHash, windowHash int
	
	h := 1
    for i := 0; i < len(pattern) - 1; i++ {
        h = (h * base) % prime
    }
    
    for i := 0; i < len(pattern); i++ {
        patternHash = (base*patternHash + int(pattern[i])) % prime
        windowHash = (base*windowHash + int(text[i])) % prime
    }
    
    var match bool
    for i := 0; i <= len(text) - len(pattern); i++ {
        if patternHash == windowHash {
            match = true
            for j := 0; j < len(pattern); j++ {
                if text[i+j] != pattern[j] {
                    match = false
                    break
                }
            }
            if match {
                result = append(result, i)
            }
        }
        
        if i < len(text) - len(pattern) {
            windowHash = (base*(windowHash-int(text[i])*h) + int(text[i+len(pattern)])) % prime
            
            if windowHash < 0 {
                windowHash += prime
            }
        }
    }
    
    return result
}
