package main

import (
	"fmt"
)

func main() {
	pattern := "AAACAAAA"

	lps := generateLPS(pattern)
	fmt.Println(lps)

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
	if len(text) == 0 || len(pattern) == 0 {
		return []int{}
	}

	patternSize := len(pattern)
	occurrenceIndexes := []int{}

	for i := 0; i <= len(text)-patternSize; i++ {
		if text[i:i+patternSize] == pattern {
			occurrenceIndexes = append(occurrenceIndexes, i)
		}
	}

	return occurrenceIndexes
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	if len(text) == 0 || len(pattern) == 0 {
		return []int{}
	}

	patternSize := len(pattern)
	occurrenceIndexes := []int{}
	lps := generateLPS(pattern)

	textIndex, patternIndex := 0, 0

	for textIndex < len(text) {
		if pattern[patternIndex] == text[textIndex] {
			textIndex++
			patternIndex++
		}

		if patternIndex == patternSize {
			occurrenceIndexes = append(occurrenceIndexes, textIndex-patternIndex)
			patternIndex = lps[patternIndex-1]
		} else if textIndex < len(text) && text[textIndex] != pattern[patternIndex] {
			if patternIndex != 0 {
				patternIndex = lps[patternIndex-1]
			} else {
				textIndex++
			}
		}
	}

	return occurrenceIndexes
}

func generateLPS(pattern string) []int {
	n := len(pattern)
	lps := make([]int, n)
	length := 0
	i := 1

	for i < n {
		if pattern[i] == pattern[length] {
			length++
			lps[i] = length
			i++
		} else {
			if length != 0 {
				length = lps[length-1] // bir önceki uzun pref. suffix'e geri dön
			} else {
				lps[i] = 0
				i++
			}
		}
	}

	return lps
}

const (
	d = 256 // Alfabe boyutu (ASCII)
	q = 101 // Büyük asal sayı (gerçek uygulamada daha büyük olmalı)
)

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	textLength := len(text)
	patternLength := len(pattern)

	if patternLength > textLength || patternLength == 0 {
		return []int{}
	}

	// h = d^(m-1) mod q (en yüksek basamağın katsayısı)
	// Overflow riskinden dolayi for dongusuyle yapiyoruz.
	// Her adimda mod aliyoruz
	h := 1
	for i := 0; i < patternLength-1; i++ {
		h = (h * d) % q
	}

	patternHash := 0
	textHash := 0

	for i := 0; i < len(pattern); i++ {
		patternHash = (d*patternHash + int(pattern[i])) % q
		textHash = (d*textHash + int(text[i])) % q
	}

	occurrenceIndexes := make([]int, 0)

	for i := 0; i <= textLength-patternLength; i++ {
		if patternHash == textHash {
			if text[i:i+patternLength] == pattern {
				occurrenceIndexes = append(occurrenceIndexes, i)
			}
		}
		if i < textLength-patternLength {
			textHash = (d*(textHash-int(text[i])*h) + int(text[i+patternLength])) % q
			if textHash < 0 {
				textHash += q
			}
		}
	}

	return occurrenceIndexes
}
