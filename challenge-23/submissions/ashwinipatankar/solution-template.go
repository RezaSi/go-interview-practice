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
	// TODO: Implement this function
	if !validate(text, pattern) {
		return []int{}
	}

	matchedIndexes := []int{}

	for i := 0; i < len(text); i++ {
		matched := true
		for j := 0; j < len(pattern); j++ {
			if i+j >= len(text) || pattern[j] != text[i+j] {
				matched = false
				break
			}
		}

		if matched {
			matchedIndexes = append(matchedIndexes, i)
		}
	}

	return matchedIndexes
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	// TODO: Implement this function
	if !validate(text, pattern) {
		return []int{}
	}

	lps := getLPSArray(pattern) //LPS ← ComputeLPS(Pattern) {build LPS table function}

	i := 0            // i ← 0
	j := 0            //j ← 0
	n := len(text)    //n ← string length
	m := len(pattern) //m ← pattern length
	results := []int{}
	for i < n {
		if pattern[j] == text[i] {
			i++
			j++

			if j == m {
				results = append(results, i-j)
				j = lps[j-1]
				continue
			}
		} else if i < n && pattern[j] != text[i] {
			if j > 0 {
				j = lps[j-1]
			} else {
				i++
			}
		}
	}
	
	return results
}

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	// TODO: Implement this function
	if !validate(text, pattern) {
		return []int{}
	}

	results := []int{}

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

	for i := 0; i <= n-m; i++ {
		isMatch := true
		if windowHash == patternHash {
			for j := 0; j < m; j++ {
				if text[i+j] != pattern[j] {
					isMatch = false
					break
				}
			}
		} else {
		    isMatch = false
		}

		if isMatch {
			results = append(results, i)
		}

		if i < n-m {
			windowHash = (base*(windowHash-int(text[i])*h) + int(text[i+m])) % prime

			// Ensure we only have positive hash values
			if windowHash < 0 {
				windowHash += prime
			}
		}
	}

	return results
}


func getLPSArray(pattern string) []int {
	lps := make([]int, len(pattern))

	//LPS ← array [size = pattern length]
	lps[0] = 0
	//LPS[0] ← 0  {LPS value of the first element is always 0}
	length := 0       //← 0  {length of previous longest proper prefix that is also a suffix}
	i := 1            // i ← 1
	m := len(pattern) // m ← length of pattern

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

func validate(text, pattern string) bool {
	if len(text) == 0 || len(pattern) == 0 || len(text) < len(pattern) {
		return false
	}

	return true
}
