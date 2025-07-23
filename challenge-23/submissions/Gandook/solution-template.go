package main

import (
	"fmt"
)

const (
	MOD  int64 = 1000000007
	BASE int64 = 257
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

func naiveCheck(text, pattern string, ind int) bool {
	for i := ind; i < ind+len(pattern); i++ {
		if text[i] != pattern[i-ind] {
			return false
		}
	}
	return true
}

// NaivePatternMatch performs a brute force search for pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func NaivePatternMatch(text, pattern string) []int {
	if len(pattern) == 0 {
		return []int{}
	}

	result := make([]int, 0)
	for i := 0; i <= len(text)-len(pattern); i++ {
		if naiveCheck(text, pattern, i) {
			result = append(result, i)
		}
	}
	return result
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	if len(pattern) == 0 || len(text) == 0 {
		return []int{}
	}

	result := make([]int, 0)
	patternKMP := make([]int, len(pattern))
	textKMP := make([]int, len(text))

	ind := -1
	patternKMP[0] = -1
	for i := 1; i < len(pattern); i++ {
		for {
			if pattern[i] == pattern[ind+1] {
				ind++
				patternKMP[i] = ind
				break
			} else if ind == -1 {
				patternKMP[i] = -1
				break
			} else {
				ind = patternKMP[ind]
			}
		}
	}

	ind = -1
	for i := 0; i < len(text); i++ {
		for {
			if text[i] == pattern[ind+1] {
				ind++
				textKMP[i] = ind
				if ind == len(pattern)-1 {
					result = append(result, i-len(pattern)+1)
					ind = patternKMP[ind]
				}
				break
			} else if ind == -1 {
				textKMP[i] = -1
				break
			} else {
				ind = patternKMP[ind]
			}
		}
	}

	return result
}

func pow(a, b int64) int64 {
	a %= MOD
	b %= MOD - 1

	if b == 0 {
		return 1
	} else if b == 1 {
		return a
	}

	num := pow(a, b/2)
	numSquared := (num * num) % MOD
	if b&1 == 0 {
		return numSquared
	} else {
		return (numSquared * a) % MOD
	}
}

func getRollingHash(s string) []int64 {
	hash := make([]int64, len(s)+1)

	for i, r := range s {
		hash[i+1] = (hash[i]*BASE + int64(r)) % MOD
	}

	return hash
}

func calcHash(rHash *[]int64, start, length int) int64 {
	return ((*rHash)[start+length] - ((*rHash)[start]*pow(BASE, int64(length)))%MOD + MOD) % MOD
}

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	if len(text) == 0 || len(pattern) == 0 {
		return []int{}
	}

	textHash := getRollingHash(text)
	patternHash := getRollingHash(pattern)
	fullPatternHash := patternHash[len(pattern)]
	result := make([]int, 0)

	for i := 0; i <= len(text)-len(pattern); i++ {
		if calcHash(&textHash, i, len(pattern)) == fullPatternHash {
			result = append(result, i)
		}
	}

	return result
}
