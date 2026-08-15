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
	matches := []int{}
	if len(pattern) == 0 || len(text) == 0 {
		return matches
	}

	for i := 0; i <= len(text)-len(pattern); i++ {
		matched := true
		for j := 0; j < len(pattern); j++ {
			if text[i+j] != pattern[j] {
				matched = false
				break
			}
		}
		if matched {
			matches = append(matches, i)
		}
	}
	return matches
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	matches := []int{}
	if len(pattern) == 0 || len(text) < len(pattern) {
		return matches
	}

	lps := computeLSP(pattern)

	i, j := 0, 0

	for i < len(text) {
		if pattern[j] == text[i] {
			i++
			j++
		}

		if j == len(pattern) {
			matches = append(matches, i-j)
			j = lps[j-1]
		}
		if i >= len(text) {
			break
		}
		if pattern[j] != text[i] {
			if j != 0 {
				j = lps[j-1]
			} else {
				i++
			}
		}
	}

	return matches
}

func computeLSP(pattern string) []int {
	l := len(pattern)
	lps := make([]int, l)
	i, j := 1, 0
	for i < l {
		if pattern[i] == pattern[j] {
			j++
			lps[i] = j
			i++
		} else if j > 0 {
			j = lps[j-1]
		} else {
			lps[i] = 0
			i++
		}
	}

	return lps
}

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	matches := []int{}
	if len(pattern) == 0 || len(text) < len(pattern) {
		return matches
	}
	m := len(pattern)
	prime := 101
	base := 256
	patternHash := 0
	windowHash := 0
	h := 1
	for i := 0; i < m-1; i++ {
		h = (h * base) % prime
	}
	for i := 0; i < m; i++ {
		patternHash = (base*patternHash + int(pattern[i])) % prime
		windowHash = (base*windowHash + int(text[i])) % prime
	}
	for i := 0; i <= len(text)-m; i++ {
		if patternHash == windowHash {
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
		if i < len(text)-m {
			windowHash = (base*(windowHash-int(text[i])*h) + int(text[i+m])) % prime
			if windowHash < 0 {
				windowHash += prime
			}
		}
	}
	return matches
}
