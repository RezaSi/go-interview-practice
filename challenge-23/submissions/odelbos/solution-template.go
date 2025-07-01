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
	result := []int{}
	tLen, pLen := len(text), len(pattern)
	if tLen == 0 || pLen == 0 || tLen < pLen {
		return result
	}

	for i := range(tLen - pLen + 1) {
		if pattern == text[i: i + pLen] {
			result = append(result, i)
		}
	}
	return result
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	result := []int{}
	tLen, pLen := len(text), len(pattern)
	if tLen == 0 || pLen == 0 || tLen < pLen {
		return result
	}

	lps := longestProperPrefix(pattern)
	i, j := 0, 0
	for i < tLen {
		if pattern[j] == text[i] {
			i++
			j++
		}
		if j == pLen {
			result = append(result, i - j)
			j = lps[j-1]
		} else if i < tLen && pattern[j] != text[i] {
			if j != 0 {
				j = lps[j - 1]
			} else {
				i++
			}
		}
	}
	return result
}

func longestProperPrefix(pattern string) []int {
	pLen := len(pattern)
	lps := make([]int, pLen)
	length, i := 0, 1
	for i < pLen {
		if pattern[i] == pattern[length] {
			length++
			lps[i] = length
			i++
		} else {
			if length != 0 {
				length = lps[length - 1]
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
	result := []int{}
	tLen, pLen := len(text), len(pattern)
	if tLen == 0 || pLen == 0 || tLen < pLen {
		return result
	}

	const d = 256
	const q = 101
	patternHash, windowHash, h := 0, 0, 1

	for i := range(pLen) {
		patternHash = (d * patternHash + int(pattern[i])) % q
		windowHash = (d * windowHash + int(text[i])) % q
	}

	for range(pLen - 1) {
		h = (h * d) % q
	}

	for i := range(tLen - pLen + 1) {
		if patternHash == windowHash {
			j := 0
			for j < pLen && text[i + j] == pattern[j] {
				j++
			}
			if j == pLen {
				result = append(result, i)
			}
		}

		if i < tLen - pLen {
			windowHash = (d * (windowHash - int(text[i]) * h) + int(text[i + pLen])) % q
			if windowHash < 0 {
				windowHash += q
			}
		}
	}
	return result
}
