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
    rText := []rune(text)
    rPattern := []rune(pattern)
    
	matches := make([]int, 0, len(rText) / 10)
	if len(rText) == 0 || len(rPattern) == 0 || len(rPattern) > len(rText) {
	    return matches
	}
	
	for i := 0; i <= len(rText)-len(rPattern); i++ {
	    j := 0
	    
	    for j < len(rPattern) && rText[i+j] == rPattern[j] {
	        j++
	    }
	    
	    if j == len(rPattern) {
	        matches = append(matches, i)
	    }
	}
	
	return matches
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
    rText := []rune(text)
    rPattern := []rune(pattern)
    
    matches := make([]int, 0, len(rText) / 10)
	if len(rText) == 0 || len(rPattern) == 0 || len(rPattern) > len(rText) {
	    return matches
	}
	
	lsp := lspBuild(rPattern)
	
	j := 0
	for i := range len(rText) {
	    for j > 0 && rText[i] != rPattern[j] {
	        j = lsp[j-1]
	    }
	    if rText[i] == rPattern[j] {
	        j++
	    }
	    if j == len(rPattern) {
	        matches = append(matches, i-len(rPattern)+1)
	        j = lsp[j-1]
	    }
	}
	return matches
}

func lspBuild(pat []rune) []int {
    lsp := make([]int, len(pat))
    
    j := 0
    for i := 1; i < len(pat); i++ {
        for j > 0 && pat[i] != pat[j] {
            j = lsp[j-1]
        }
        if pat[i] == pat[j] {
            j++
        }
        lsp[i] = j
    }
    
    return lsp
}

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	rText, rPat := []rune(text), []rune(pattern)
	nT, nP := len(rText), len(rPat)
	
	matches := make([]int, 0, nT / 10)
	if nT == 0 || nP == 0 || nT < nP {
	    return matches
	}
	
	const base = 65536
	const prime = 10000003
	
	patHash, curHash := 0, 0
	high := 1
	
	for i := 0; i < nP-1; i++ {
	    high = (high * base) % prime
	}
	
	for i := 0; i < nP; i++ {
	    patHash = (base*patHash + int(rPat[i])) % prime
	    curHash = (base*curHash + int(rText[i])) % prime
	}
	
	for i := 0; i <= nT-nP; i++ {
	    if patHash == curHash {
	        match := true
	        for j := 0; j < nP; j++ {
	            if rText[i+j] != rPat[j] {
	                match = false
	                break
	            }
	        }
	        if match {
	            matches = append(matches, i)
	        }
	    }
	    
	    if i < nT-nP {
	        curHash = base * (curHash - int(rText[i])*high)
	        curHash %= prime
	        curHash = (curHash + int(rText[i+nP])) % prime
	        if curHash < 0 {
	            curHash += prime
	        }
 	    }
	}
	
	return matches
}
