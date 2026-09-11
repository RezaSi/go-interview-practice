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

	pIndx := []int{}

	if len(text) == 0 || len(pattern) == 0 || len(pattern) > len(text) {
		return pIndx
	}

	for i := 0; i <= len(text)-len(pattern); i++ {

		j := 0

		for j < len(pattern) && text[i+j] == pattern[j] {
			j++
		}

		if j == len(pattern) {
			pIndx = append(pIndx, i)
		}
	}

	return pIndx
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {

    matches := []int{}
   if len(text) == 0 || len(pattern) == 0 || len(pattern) > len(text) {
        return matches
    }
	
	lps := func(pattern string) []int {
		// creating the longest prefix from pattern for processing string

		leng := 0
		res := make([]int, len(pattern))
		res[0] = 0
		i := 1

		for i < len(pattern) {
			if pattern[i] == pattern[leng] {
				leng++
				res[i] = leng
				i++
			} else {
				if leng != 0 {
					leng = res[leng-1]
				}
				res[i] = 0
				i++

			}

		}

		return res
	}

	lenT := len(text)
	lenP := len(pattern)

	i := 0
	j := 0

	for i < lenT {
		if text[i] == pattern[j] {
			i++
			j++

		}

		if j == lenP {
			matches = append(matches, (i - j))
			j = lps(pattern)[j-1]

		} else if i < lenT && pattern[j] != text[i] {

			if j != 0 {
				j = lps(pattern)[j-1]
			} else {
				i++
			}
		}
	}

	return matches
}

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	matches := []int{}

	// Handle edge cases
	if len(pattern) == 0 || len(text) < len(pattern) {
		return matches
	}

	n := len(text)
	m := len(pattern)

	// Large prime number to avoid hash collisions
	prime := 101

	// Base value for the hash function
	base := 256

	// Hash value for pattern and initial window
	patternHash := 0
	windowHash := 0

	// Highest power of base that we need
	h := 1
	for i := 0; i < m-1; i++ {
		h = (h * base) % prime
	}

	// Calculate initial hash values
	for i := 0; i < m; i++ {
		patternHash = (base*patternHash + int(pattern[i])) % prime
		windowHash = (base*windowHash + int(text[i])) % prime
	}

	// Slide the pattern over text one by one
	for i := 0; i <= n-m; i++ {
		// Check if hash values match
		if patternHash == windowHash {
			// Verify the match character by character
			match := true
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

		// Calculate hash value for next window
		if i < n-m {
			windowHash = (base*(windowHash-int(text[i])*h) + int(text[i+m])) % prime

			// Ensure we only have positive hash values
			if windowHash < 0 {
				windowHash += prime
			}
		}
	}

	return matches
}
