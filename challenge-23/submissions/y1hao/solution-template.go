package main

import (
	"fmt"
	"strings"
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
	if text == "" || pattern == "" {
		return []int{}
	}

	result := []int{}
	for i := 0; i <= len(text)-len(pattern); i++ {
		if strings.HasPrefix(text[i:], pattern) {
			result = append(result, i)
		}
	}
	return result
}

// I gave up on KMP and the followings are copied from learning :)

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	if text == "" || pattern == "" {
		return []int{}
	}

	n := len(text)
	m := len(pattern)

	// Preprocess the pattern
	lps := computeLPSArray(pattern)

	i := 0 // Index for text
	j := 0 // Index for pattern

	matches := []int{}
	for i < n {
		// Current characters match, move both pointers forward
		if pattern[j] == text[i] {
			i++
			j++
		}

		// Found a complete match
		if j == m {
			matches = append(matches, i-j)
			// Use lps to shift pattern for next match
			j = lps[j-1]
		} else if i < n && pattern[j] != text[i] {
			// Mismatch after j matches
			if j != 0 {
				// Use lps to shift pattern
				j = lps[j-1]
			} else {
				// No match found, move to next character in text
				i++
			}
		}
	}

	return matches
}

func computeLPSArray(pattern string) []int {
	m := len(pattern)
	lps := make([]int, m)

	// Length of the previous longest prefix suffix
	length := 0
	i := 1

	// The loop calculates lps[i] for i = 1 to m-1
	for i < m {
		if pattern[i] == pattern[length] {
			length++
			lps[i] = length
			i++
		} else {
			// This is the tricky part
			if length != 0 {
				length = lps[length-1]
				// Note: We do not increment i here
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
	if text == "" || pattern == "" || len(text) < len(pattern) {
		return []int{}
	}

	factor := 1
	for range len(pattern) - 1 {
		factor = (factor * 256) % 101
	}

	target := hash(pattern)
	current := hash(text[:len(pattern)])
	results := []int{}
	if target == current {
		results = append(results, 0)
	}
	for i := 1; i <= len(text)-len(pattern); i++ {
		current = nextHash(current, factor, text[i-1], text[i+len(pattern)-1])
		if current == target {
			if strings.HasPrefix(text[i:], pattern) {
				results = append(results, i)
			}
		}
	}

	return results
}

func hash(text string) int {
	h := 0
	for _, b := range text {
		h = (h*256 + int(b)) % 101
	}
	return h
}

func nextHash(h int, factor int, first, current byte) int {
	h -= int(first) * factor
	h = (h*256 + int(current)) % 101
	if h < 0 {
		h += 101
	}
	return h
}
