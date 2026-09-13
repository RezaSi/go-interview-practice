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

	pl := len(pattern)
	tl := len(text)

	res := []int{}

	if pl == 0 || tl == 0 || pl > tl {
		return res
	}

	for i := 0; i < tl-pl+1; i++ {
		match := true
		for j := 0; j < pl; j++ {
			if text[i+j] != pattern[j] {
				match = false
				break
			}
		}
		if match {
			res = append(res, i)
		}
	}

	return res
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.

func prefix(str string) []int {
	length := len(str)
	res := make([]int, length, length)
	k := 0

	for i := 1; i < length; i++ {
		for k != 0 && str[i] != str[k] {
			k = res[k-1]
		}
		if str[i] == str[k] {
			k += 1
		}
		res[i] = k
	}

	return res
}

func KMPSearch(text, pattern string) []int {

	pl := len(pattern)
	tl := len(text)

	pif := prefix(pattern)

	res := []int{}

	if pl == 0 || tl == 0 || pl > tl {
		return res
	}

	k := 0

	for i := 0; i < tl; i++ {

		for k != 0 && pattern[k] != text[i] {
			k = pif[k-1]
		}
		if pattern[k] == text[i] {
			k++
		}
		if k == pl {
			res = append(res, i-pl+1)
			k = pif[pl-1]
		}

	}

	return res
}

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	const (
		B = 256
		M = 1000000007
	)
	pl := len(pattern)
	tl := len(text)

	res := []int{}

	if pl == 0 || tl == 0 || pl > tl {
		return res
	}

	var highPower int64 = 1
	for i := 0; i < pl-1; i++ {
		highPower = (highPower * B) % M
	}

	var patternHash, textHash int64

	for i := 0; i < pl; i++ {
		patternHash = (patternHash*B + int64(pattern[i])) % M
		textHash = (textHash*B + int64(text[i])) % M
	}

	for i := 0; i <= tl-pl; i++ {
		if patternHash == textHash {
			match := true
			for j := 0; j < pl; j++ {
				if text[i+j] != pattern[j] {
					match = false
					break
				}
			}
			if match {
				res = append(res, i)
			}

		}
		if i < tl-pl {
			firstChar := int64(text[i])
			nextChar := int64(text[i+pl])

			textHash = ((textHash-firstChar*highPower%M)*B + nextChar) % M
			if textHash < 0 {
				textHash += M
			}

		}
	}

	return res
}
