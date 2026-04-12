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
	result := []int{} // Initialize as empty slice, not nil
	n := len(text)
	m := len(pattern)

	if m == 0 || n < m {
		return result
	}

	for i := 0; i <= n-m; i++ {
		j := 0
		for j < m && text[i+j] == pattern[j] {
			j++
		}
		if j == m {
			result = append(result, i)
		}
	}
	return result
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	result := []int{} // Initialize as empty slice, not nil
	n := len(text)
	m := len(pattern)

	if m == 0 || n < m {
		return result
	}

	// Build the prefix table (also called failure function or lps)
	lps := computeLPS(pattern)

	i := 0 // index for text
	j := 0 // index for pattern

	for i < n {
		if pattern[j] == text[i] {
			i++
			j++
		}

		if j == m {
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

// computeLPS computes the Longest Prefix Suffix array for KMP algorithm
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
	result := []int{} // Initialize as empty slice, not nil
	n := len(text)
	m := len(pattern)

	if m == 0 || n < m {
		return result
	}

	// Prime number for hashing
	prime := 101
	// Base for the polynomial rolling hash
	d := 256

	// Calculate hash values
	patternHash := 0
	textHash := 0
	h := 1

	// The value of h would be "pow(d, m-1) % prime"
	for i := 0; i < m-1; i++ {
		h = (h * d) % prime
	}

	// Calculate the hash value of pattern and first window of text
	for i := 0; i < m; i++ {
		patternHash = (d*patternHash + int(pattern[i])) % prime
		textHash = (d*textHash + int(text[i])) % prime
	}

	// Slide the pattern over text one by one
	for i := 0; i <= n-m; i++ {
		// Check the hash values of current window of text and pattern
		if patternHash == textHash {
			// Check for characters one by one if hash matches
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
			textHash = (d*(textHash-int(text[i])*h) + int(text[i+m])) % prime

			// We might get negative value of textHash, converting it to positive
			if textHash < 0 {
				textHash += prime
			}
		}
	}
	return result
}