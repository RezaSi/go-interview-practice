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
	// TODO: Implement this function
	textLen := len(text)
	patternLen := len(pattern)
	result := []int{}
	if patternLen == 0 {
	    return result
	}

	for i := 0; i <= (textLen - patternLen); i++ {
	    if text[i] == pattern[0]{
	        for j:= 0; j < patternLen; j++ {
	            if text[i+j] != pattern[j] {
	                break
	            }
	            if j == patternLen -1 {
	                result = append(result, i)
	            }
	        }
	    }   
	}
	return result
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	// TODO: Implement this function
	textLen := len(text)
	patternLen := len(pattern)
	result := []int{}
	if patternLen == 0 {
	    return result
	}
	table := KMPTable(pattern)

	for i, j := 0, 0; i < textLen; {
	    if pattern[j] == text[i] {
	        i, j = i+1, j+1
	        if j == patternLen {
	            result = append(result, i - j)
	            j = table[j]
	        }
	    } else {
	        j = table[j]
	        if j < 0 {
	            i, j = i+1, j+1
	        }
	    }
	}
	return result
}

func KMPTable(pattern string) []int {
    patternLen := len(pattern)
    if patternLen == 0 {
        return []int{}
    }
    table := make([]int, patternLen + 1)
    table[0] = -1
    cnd := 0
    for i:= 1; i < patternLen; i, cnd = i+1, cnd+1 {
        if pattern[i] == pattern[cnd] {
            table[i] = table[cnd]
        } else {
            table[i] = cnd
            for cnd >= 0 && pattern[i] != pattern[cnd] {
                cnd = table[cnd]
            }
        }
    }
    table[patternLen] = cnd
    return table
}

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	// TODO: Implement this function
	textLen := len(text)
	patternLen := len(pattern)
	result := []int{}
	if patternLen == 0 || textLen == 0 || textLen < patternLen {
	    return result
	}
	
	base, prime := 256, 101
	highestBasePower := 1
	for idx := 0; idx < patternLen - 1; idx++ {
	    highestBasePower = (highestBasePower * base) % prime
	}
	
	patternHash, sliceHash := 0, 0
	for idx := range pattern {
	    patternHash = (patternHash * base + int(pattern[idx])) % prime
	    sliceHash = (sliceHash * base + int(text[idx])) % prime
	}
	
	for idx := 0; idx <= (textLen - patternLen); idx++ {
	    if sliceHash == patternHash && text[idx:idx+patternLen] == pattern {
	        result = append(result, idx)
	    }
	    if idx < (textLen - patternLen){
	        sliceHash = ((sliceHash - int(text[idx]) * highestBasePower) * base + int(text[idx + patternLen])) % prime
	        if sliceHash < 0 {
	            sliceHash += prime
	        }
	    }
    }
	return result
}
