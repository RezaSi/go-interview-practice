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
	ans := []int{}
	if len(text) == 0 || len(pattern) == 0 {
		return ans
	}

	patternLen := len(pattern)
	textLen := len(text)

	for i := 0; i <= textLen-patternLen; i++ {
		j := 0
		for j = 0; j < patternLen; j++ {
			if text[i+j] != pattern[j] {
				break
			}
		}

		if j == patternLen {
			ans = append(ans, i)
		}
	}

	return ans
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	patternLen := len(pattern)
	textLen := len(text)
	ans := []int{}

	if patternLen > textLen || patternLen == 0 {
		return ans
	}

	lsp := LPS(pattern)
	i, j := 0, 0

	for i < textLen {
		if text[i] == pattern[j] {
			i += 1
			j += 1

			if j == patternLen {
				ans = append(ans, i-j)
				j = lsp[j-1]
			}
		} else {
			if j != 0 {
				j = lsp[j-1]
			} else {
				i += 1
			}
		}
	}

	return ans
}

func LPS(pattern string) []int {
	patternLen := len(pattern)
	lsp := make([]int, patternLen)

	i := 1
	j := 0

	// ABABCABAB
	for i < patternLen {
		if pattern[i] == pattern[j] {
			j += 1
			lsp[i] = j
			i += 1
		} else {
			if j != 0 {
				j = lsp[j-1]
			} else {
				i += 1
			}
		}
	}

	return lsp
}

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	n, m := len(text), len(pattern)
	matches := []int{}

	if m == 0 || n == 0 || m > n {
		return matches
	}

	base := 256
	mod := 1_000_000_007

	h := 1
	for range m - 1 {
		h = (h * base) % mod
	}

	patternHash, windowHash := 0, 0

	for i := range m {
		patternHash = (patternHash*base + int(pattern[i])) % mod
		windowHash = (windowHash*base + int(text[i])) % mod
	}

	for i := 0; i <= n-m; i++ {
		if patternHash == windowHash {
			if text[i:i+m] == pattern {
				matches = append(matches, i)
			}
		}

		if i < n-m {
			windowHash = (windowHash - int(text[i])*h%mod + mod) % mod
			windowHash = (windowHash*base + int(text[i+m])) % mod
		}
	}

	return matches
}
