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

	n := len(text)
	m := len(pattern)
	result := make([]int, 0)

	if n < m || n <= 0 || m <= 0 {
		return result
	}

	for i := 0; i <= n-m; i++ {
		find := true
		for j := 0; j < m; j++ {
			if text[i+j] != pattern[j] {
				find = false
				break
			}
		}
		if find {
			result = append(result, i)
		}

	}

	return result
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	// TODO: Implement this function

	n := len(text)
	m := len(pattern)
	lps := computeLPS(pattern)
	result := make([]int, 0)
	i, j := 0, 0

	if n < m || n <= 0 || m <= 0 {
		return result
	}

	for i < n {
		if text[i] == pattern[j] {
			i++
			j++
			if j == m {
				result = append(result, i-j)
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
	return result

}

func computeLPS(pattern string) []int {
	m := len(pattern)
	lps := make([]int, m)
	length := 0
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
	// TODO: Implement this function

	d := 256 // Размер алфавита
	q := 101 // Простое число для хэширования

	n := len(text)
	m := len(pattern)
	result := make([]int, 0)

	if n < m || n <= 0 || m <= 0 {
		return result
	}

	// Вычисляем h = d^(m-1) mod q
	h := 1
	for i := 0; i < m-1; i++ {
		h = (h * d) % q
	}

	// Вычисляем хэш паттерна и первого окна в тексте
	pHash, tHash := 0, 0
	for i := 0; i < m; i++ {
		pHash = (d*pHash + int(pattern[i])) % q
		tHash = (d*tHash + int(text[i])) % q
	}

	// Скользящее окно по тексту
	for i := 0; i <= n-m; i++ {
		// Если хэши совпали, проверяем символы
		if pHash == tHash {
			find := true
			for j := 0; j < m; j++ {
				if text[i+j] != pattern[j] {
					find = false
					break
				}
			}
			if find {
				result = append(result, i)
			}
		}

		// Пересчитываем хэш для следующего окна
		if i < n-m {
			tHash = (d*(tHash-int(text[i])*h) + int(text[i+m])) % q
			if tHash < 0 {
				tHash += q
			}
		}
	}
	return result
}
