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

// Time complexity: O(nm) (n - text, m - pattern)
// Space complexity: O(k) (k - number of matches)

func NaivePatternMatch(text, pattern string) []int {
	matches := []int{}
	
	// edge cases
	if len(pattern) == 0 || len(text) < len(pattern) {
	    return matches
	}
	
	for i := 0; i <= len(text) - len(pattern); i++ {
	    j := 0
	    
	    for j < len(pattern) && text[i + j] == pattern[j] {
	        j++
	    }
	    
	    // if j is the end of the pattern, match is found
	    if j == len(pattern) {
	        matches = append(matches, i)
	    }
	}
	
	return matches
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.

// Time complexity: O(n + m) (n - text, m - pattern)
// Space complexity: O(m), O(k) (m for LPS array, k - number of matches)
func KMPSearch(text, pattern string) []int {
	matches := []int{}
	
	// edge cases
	if len(pattern) == 0 || len(text) < len(pattern) {
	    return matches
	}
	
	n := len(text)
	m := len(pattern)
	
	lps := computeLPSArray(pattern)
	
	i := 0 // text
	j := 0 // pattern
	
	for i < n {
	    if pattern[j] == text[i] {
	        // match, move both pointers forward
	        i++
	        j++
	    }
	    
	    if j == m {
	        // complete match
	        matches = append(matches, i-j)
	        // shift pattern for next match
	        j = lps[j-1]
	    } else if i < n && pattern[j] != text[i] {
	        // mismatch after j matches
	        if j != 0 {
	            j = lps[j-1]
	        } else {
	            // no match, to next char in text
	            i++
	        }
	    }
	}
	
	return matches
}

// KMP helper, longest prefix suffix (LPS) array, avoids redundant comparisons
func computeLPSArray(pattern string) []int {
    m := len(pattern)
    lps := make([]int, m) // slice, length is the lenght of the pattern
    
    // length of the prev lps
    length := 0
    i := 1
    
    // calc lps[i] for i = 1 to m-1
    for i < m {
        if pattern[i] == pattern[length] {
            length++
            lps[i] = length
            i++
        } else {
            if length != 0 {
                length = lps[length-1]
                // do not incr i here
            } else {
                lps[i] = 0
                i++
            }
        }
    }
    
    return lps
}

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.


// Time complexity average: O(n + m) (n - text, m - pattern)
// Time complexity worst (too many hash collisions): O(nm) (n - text, m - pattern)
// Space complexity: O(k) (k - number of matches)

func RabinKarpSearch(text, pattern string) []int {
	matches := []int{}
	
	// edge cases
	if len(pattern) == 0 || len(text) < len(pattern) {
	    return matches
	}
	
	n := len(text)
	m := len(pattern)
	
	// large prime to avoid hash collisions
	prime := 101
	
	// for hash fn
	base := 256
	
	patternHash := 0
	windowHash := 0
	
	// highest power of base needed
	h := 1
	for i := 0; i < m-1; i++ {
	    h = (h * base) % prime
	}
	
	// initial hash values
	for i := 0; i < m; i++ {
	    patternHash = (base*patternHash + int(pattern[i])) % prime
	    windowHash = (base*windowHash + int(text[i])) % prime
	}
	
	for i := 0; i <= n-m; i++ {
	    if patternHash == windowHash {
	        match := true
	        // verify
	        for j := 0; j < m; j++ {
	            if text[i+j] != pattern[j] {
	                match = false
	                break
	            }
	        }
	        if match {
	            matches = append(matches, i)
	        }
	    }
	    // hash for next window
	    if i < n-m {
	        windowHash = (base*(windowHash-int(text[i])*h) + int(text[i+m])) % prime
	        
	        if windowHash < 0 {
	            windowHash += prime
	        }
	    }
	}
	
	
	return matches
}
