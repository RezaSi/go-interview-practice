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
	if len(pattern) == 0 {
		return result
	}
	for i := 0; i+len(pattern)-1 < len(text); i++ {
		if text[i:i+len(pattern)] == pattern {
			result = append(result, i)
		}
	}
	return result
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	result := []int{}
	if len(pattern) == 0 {
		return result
	}
	fullText := pattern + "#" + text
	n := len(fullText)
	pi := make([]int, n, n)
	for i := 1; i < n; i++ {
		j := pi[i-1]
		for j > 0 && fullText[j] != fullText[i] {
			j = pi[j-1]
		}
		if fullText[j] == fullText[i] {
			j++
		}
		pi[i] = j
	}
	for i := len(pattern) + 1; i < n; i++ {
		if pi[i] == len(pattern) {
			result = append(result, i-2*len(pattern))
		}
	}
	return result
}

const p int64 = 31
const mod = int64(1e7 + 7)

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	result := []int{}
	if len(pattern) == 0 || len(pattern) > len(text) {
		return result
	}
	n := len(text)
	pows := make([]int64, n, n)
	pows[0] = 1
	for i := 1; i < n; i++ {
		pows[i] = (pows[i-1] * p) % mod
	}
	ht := make([]int64, n, n)
	for i, ch := range text {
		ht[i] = (int64(ch) * pows[i]) % mod
		if i > 0 {
			ht[i] += ht[i-1]
		}
	}
	var hp int64 = 0
	for i, ch := range pattern {
		hp += (int64(ch) * pows[i]) % mod
	}
	for i := 0; i+len(pattern)-1 < n; i++ {
		curr_ht := ht[i+len(pattern)-1]
		if i > 0 {
			curr_ht -= ht[i-1]
		}
		curr_ht = ((curr_ht % mod) + mod) % mod
		if curr_ht == (hp*pows[i])%mod && text[i:i+len(pattern)] == pattern {
			result = append(result, i)
		}
	}
	return result
}
