package main

import (
	"fmt"
)

// NaivePatternMatch performs a brute force search for pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func NaivePatternMatch(text, pattern string) []int {
    n := len(text)
    m := len(pattern)
    
    // Handle edge cases - always return non-nil slice
    if m == 0 || n < m {
        return []int{}
    }
    
    result := []int{}
    
    // Check each possible starting position
    for i := 0; i <= n-m; i++ {
        j := 0
        // Check if pattern matches starting at position i
        for j < m && text[i+j] == pattern[j] {
            j++
        }
        // If we reached the end of pattern, we found a match
        if j == m {
            result = append(result, i)
        }
    }
    
    return result
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
    n := len(text)
    m := len(pattern)
    
    // Handle edge cases - always return non-nil slice
    if m == 0 || n < m {
        return []int{}
    }
    
    // Preprocess the pattern to create the prefix table (lps)
    lps := computeLPS(pattern)
    
    result := []int{}
    i, j := 0, 0 // i for text, j for pattern
    
    for i < n {
        if pattern[j] == text[i] {
            i++
            j++
        }
        
        if j == m {
            // Pattern found at index i-j
            result = append(result, i-j)
            j = lps[j-1]
        } else if i < n && pattern[j] != text[i] {
            if j != 0 {
                j = lps[j-1]
            } else {
                i++
            }
        }
    }
    
    return result
}

// computeLPS computes the longest prefix suffix table for KMP algorithm
func computeLPS(pattern string) []int {
    m := len(pattern)
    lps := make([]int, m)
    length := 0 // length of the previous longest prefix suffix
    i := 1
    
    for i < m {
        if pattern[i] == pattern[length] {
            length++
            lps[i] = length
            i++
        } else {
            if length != 0 {
                length = lps[length-1]
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
func RabinKarpSearch(text, pattern string) []int {
    n := len(text)
    m := len(pattern)
    
    // Handle edge cases - always return non-nil slice
    if m == 0 || n < m {
        return []int{}
    }
    
    result := []int{}
    
    // Prime number for hashing
    prime := 101
    // Base for the polynomial rolling hash (number of characters in alphabet)
    base := 256
    
    // Calculate hash value for pattern and first window of text
    patternHash := 0
    textHash := 0
    h := 1
    
    // The value of h would be "base^(m-1) % prime"
    for i := 0; i < m-1; i++ {
        h = (h * base) % prime
    }
    
    // Calculate initial hash values
    for i := 0; i < m; i++ {
        patternHash = (base*patternHash + int(pattern[i])) % prime
        textHash = (base*textHash + int(text[i])) % prime
    }
    
    // Slide the pattern over text one by one
    for i := 0; i <= n-m; i++ {
        // Check the hash values first
        if patternHash == textHash {
            // If hash matches, check characters one by one
            match := true
            for j := 0; j < m; j++ {
                if text[i+j] != pattern[j] {
                    match = false
                    break
                }
            }
            if match {
                result = append(result, i)
            }
        }
        
        // Calculate hash value for next window of text
        if i < n-m {
            textHash = (base*(textHash - int(text[i])*h) + int(text[i+m])) % prime
            
            // We might get negative value of textHash, converting it to positive
            if textHash < 0 {
                textHash += prime
            }
        }
    }
    
    return result
}

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