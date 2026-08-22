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
	res := make([]int, 0)

	if pattern == "" {
		return res
	}
	for i := 0; i <= len(text)-len(pattern); i++ {
		if text[i:i+len(pattern)] == pattern {
			res = append(res, i)
		}
	}
	return res
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	// TODO: Implement this function
	if pattern == "" {
		return []int{}
	}
	var n, m = len(text), len(pattern)
	lps := make([]int, m)
	res := make([]int, 0)
	constructLps(pattern, lps)

	var i, j = 0, 0

	for i < n {
		if text[i] == pattern[j] {
			i++
			j++

			if j == m {
				res = append(res, i-j)
				j = lps[j-1]
			}
		} else {
			if j != 0 {
				j = lps[j-1]
			} else {
				i++
			}
		}
	}

	return res
}

func constructLps(pattern string, lps []int) {
	var i, length = 1, 0

	for i < len(pattern) {
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
}

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	res := []int{}
	n, m := len(text), len(pattern)
	if m == 0 || m > n {
		return res
	}

	const d = 256
	const q = 101

	// h = d^(m-1) % q
	h := 1
	for i := 0; i < m-1; i++ {
		h = (h * d) % q
	}

	patHash, txtHash := 0, 0
	for i := 0; i < m; i++ {
		patHash = (d*patHash + int(pattern[i])) % q
		txtHash = (d*txtHash + int(text[i])) % q
	}

	for i := 0; i <= n-m; i++ {
		if patHash == txtHash {
			if text[i:i+m] == pattern {
				res = append(res, i)
			}
		}
		if i < n-m {
			txtHash = (d*(txtHash-int(text[i])*h) + int(text[i+m])) % q
			if txtHash < 0 {
				txtHash += q
			}
		}
	}
	return res
}
