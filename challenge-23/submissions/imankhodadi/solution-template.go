// Challenge 23: String Pattern Matching
package main

// NaivePatternMatch performs a brute force search for pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func NaivePatternMatch(text, pattern string) []int {
	matches := []int{}
	if len(pattern) == 0 || len(text) < len(pattern) {
		return matches
	}
	// Check each possible position in the text
	for i := 0; i <= len(text)-len(pattern); i++ {
		j := 0
		// Check if the pattern matches at this position
		for j < len(pattern) && text[i+j] == pattern[j] {
			j++
		}
		// If j reached the end of the pattern, we found a match
		if j == len(pattern) {
			matches = append(matches, i)
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

// For KMP: Optimize the LPS computation
func computeLPSOptimized(pattern string) []int {
	m := len(pattern)
	lps := make([]int, m)
	for i, length := 1, 0; i < m; {
		if pattern[i] == pattern[length] {
			length++
			lps[i] = length
			i++
		} else if length != 0 {
			length = lps[length-1]
		} else {
			lps[i] = 0
			i++
		}
	}
	return lps
}

/*
Preprocess the pattern to build a partial match table (also called the "LPS" or "π" table)
Use this table to determine how far to shift the pattern when a mismatch occurs
Never backtrack in the text - each character in the text is examined exactly once
Creating the LPS (Longest Prefix Suffix) Array
The LPS array helps determine the longest proper prefix of the pattern that is also a suffix of the pattern up to each position.
This information is used to avoid redundant comparisons.

Complexity of the KMP Algorithm
Time Complexity: O(n+m) where n is the length of the text and m is the length of the pattern
Space Complexity: O(m) for the LPS array plus O(k) for storing the matches
The KMP algorithm is much more efficient than the naive approach for texts with many potential matches, especially for longer patterns.*/

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.

func handleEdgeCases(text, pattern string) ([]int, bool) {
	// Empty pattern
	if len(pattern) == 0 {
		return []int{}, true
	}
	// Pattern longer than text
	if len(pattern) > len(text) {
		return []int{}, true
	}
	// Empty text but non-empty pattern
	if len(text) == 0 {
		return []int{}, true
	}
	// Continue with normal processing
	return nil, false
}

func KMPSearch(text, pattern string) []int {
	matches := []int{}

	// Handle edge cases
	if len(pattern) == 0 || len(text) < len(pattern) {
		return matches
	}

	n := len(text)
	m := len(pattern)

	// Preprocess the pattern
	lps := computeLPSArray(pattern)

	i := 0 // Index for text
	j := 0 // Index for pattern

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

/*
The Rabin-Karp algorithm uses hashing to find pattern matches more efficiently. Instead of comparing each character, it compares hash values of the pattern and substrings of the text.

How the Rabin-Karp Algorithm Works
Compute the hash value of the pattern
Compute hash values for all possible m-length substrings of the text using a rolling hash function
Compare the hash value of the pattern with the hash value of each substring
If the hash values match, verify the actual strings character by character
Implementing a Rolling Hash Function
A rolling hash function allows us to compute the hash value of the next substring in constant time by using the hash value of the current substring:

// Remove the leftmost character and add the rightmost character
newHash = (oldHash - oldChar * pow) * base + newChar
*/
// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.

// For Rabin-Karp: Use better hash function to reduce collisions
func betterHash(s string) uint64 {
	var hash uint64 = 0
	for i := 0; i < len(s); i++ {
		hash = hash*31 + uint64(s[i])
	}
	return hash
}

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
