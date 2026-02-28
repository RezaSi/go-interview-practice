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
	// Convert text and pattern to rune
	rText := []rune(text)
	rPattern := []rune(pattern)

	// Calculate rune length
	lenPattern := len(rPattern)
	lenText := len(rText)

	// Check for emty text pattern and pattern no more than text
	if text == "" || pattern == "" || lenPattern > lenText {
		return []int{}
	}

	// Init result slise
	res := []int{}

	// Loop through the text
	for i := 0; i <= lenText-lenPattern; i++ {
		match := true

		// loop through the pattern and text to compare each rune
		for j := range lenPattern {
			if rText[i+j] != rPattern[j] {
				match = false
				break
			}
		}

		// If loop ends without breaks and match it true - append index to result
		if match {
			res = append(res, i)
		}
	}
	return res
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	// Convert text and pattern to rune
	rText := []rune(text)
	rPattern := []rune(pattern)

	// Calculate rune length
	lenPattern := len(rPattern)
	lenText := len(rText)

	// Check for emty text pattern and pattern no more than text
	if text == "" || pattern == "" || lenPattern > lenText {
		return []int{}
	}

	// init lps slise
	res := []int{}
	lps := make([]int, len(pattern))
	i := 1
	l := 0

	// Loop through the pattern to calculate LPS slise
	for i < lenPattern {
		if rPattern[l] == rPattern[i] {
			l++
			lps[i] = l
			i++
		} else {
			if l != 0 {
				l = lps[l-1]
				// l--
			} else {
				lps[i] = 0
				i++
			}
		}
	}

	// Loop for KMP search
	i, j := 0, 0
	for i < lenText {
		if rText[i] == rPattern[j] {
			i++
			j++
		}
		if j == lenPattern {
			res = append(res, i-j)
			j = lps[j-1]
		} else if i < lenText && text[i] != pattern[j] {
			if j != 0 {
				j = lps[j-1]
			} else {
				i++
			}
		}
	}
	return res
}

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	const base = 256
	const mod = 101 // prime number

	// init res slise
	res := []int{}

	// Convert text and pattern to rune
	rText := []rune(text)
	rPattern := []rune(pattern)

	// Calculate rune length
	lenPattern := len(rPattern)
	lenText := len(rText)

	// Check for emty text pattern and pattern no more than text
	if text == "" || pattern == "" || lenPattern > lenText {
		return []int{}
	}

	var patternHash, textHash, h int

	// Calculate hash for pattern h = base^(m-1) % mod
	h = 1
	for i := 0; i < lenPattern-1; i++ {
		h = (h * base) % mod
	}

	// Initial hash
	for i := range lenPattern {
		patternHash = (base*patternHash + int(pattern[i])) % mod
		textHash = (base*textHash + int(text[i])) % mod
	}

	for i := 0; i <= lenText-lenPattern; i++ {
		if patternHash == textHash {
			match := true
			for j := range lenPattern {
				if text[i+j] != pattern[j] {
					match = false
					break
				}
			}
			if match {
				res = append(res, i)
			}
		}

		if i < lenText-lenPattern {
			textHash = (base*(textHash-int(text[i])*h) + int(text[i+lenPattern])) % mod
			if textHash < 0 {
				textHash += mod
			}
		}
	}

	return res
}
