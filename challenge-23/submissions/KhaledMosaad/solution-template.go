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
	if len(text) < len(pattern) || len(pattern) == 0 || len(text) == 0 {
		return []int{}
	}

	ret := make([]int, 0)
	for i := range text {
		if i+len(pattern) <= len(text) && text[i] == pattern[0] {
			exist := true
			for j := range pattern {
				if pattern[j] != text[i+j] {
					exist = false
					break
				}
			}

			if exist {
				ret = append(ret, i)
			}
		}
	}
	return ret
}

// Building longest prefix suffix slice
// Example pattern = abdabd, lps = [0,0,]
func BuildLPS(pattern string) []int {
	lps := make([]int, len(pattern))
	length, i := 0, 1
	for i < len(pattern) {
		if pattern[i] == pattern[length] { // if the current position = the prefix length
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

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	if len(text) < len(pattern) || len(pattern) == 0 || len(text) == 0 {
		return []int{}
	}
	lps := BuildLPS(pattern)
	ret := make([]int, 0)

	for i, j := 0, 0; i < len(text); {
		if text[i] == pattern[j] {
			i++
			j++
		}

		if j == len(pattern) {
			ret = append(ret, i-j)
			j = lps[j-1]
		} else if i < len(text) && text[i] != pattern[j] {
			if j != 0 {
				j = lps[j-1]
			} else {
				i++
			}
		}
	}
	return ret
}

const (
	base = 256 // base: number of characters (ASCII)
	mod  = 101 // a prime number for modulo to reduce collisions
)

// hash calculates the hash of the given string
func hash(s string, length int) int {
	h := 0
	for i := 0; i < length; i++ {
		h = (h*base + int(s[i])) % mod
	}
	return h
}

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	if len(text) < len(pattern) || len(pattern) == 0 || len(text) == 0 {
		return []int{}
	}

	m, n := len(pattern), len(text)
	patternHash := hash(pattern, m)
	textHash := hash(text, m)

	// precompute (base^(m-1)) % mod for use in rolling hash
	h := 1
	for i := 0; i < m-1; i++ {
		h = (h * base) % mod
	}

	result := make([]int, 0)
	for i := 0; i <= n-m; i++ {
		if patternHash == textHash {
			// Confirm match to avoid false positive due to hash collision
			if text[i:i+m] == pattern {
				result = append(result, i)
			}
		}

		// compute hash for next window
		if i < n-m {
			textHash = (base*(textHash-int(text[i])*h) + int(text[i+m])) % mod
			if textHash < 0 {
				textHash += mod // make sure hash is positive
			}
		}
	}
	return result
}
