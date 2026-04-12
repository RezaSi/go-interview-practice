// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	// Add any necessary imports here
	"fmt"
	"regexp"
	"strings"
)

// CountWordFrequency takes a string containing multiple words and returns
// a map where each key is a word and the value is the number of times that
// word appears in the string. The comparison is case-insensitive.
//
// Words are defined as sequences of letters and digits.
// All words are converted to lowercase before counting.
// All punctuation, spaces, and other non-alphanumeric characters are ignored.
//
// For example:
// Input: "The quick brown fox jumps over the lazy dog."
// Output: map[string]int{"the": 2, "quick": 1, "brown": 1, "fox": 1, "jumps": 1, "over": 1, "lazy": 1, "dog": 1}
func CountWordFrequency(text string) map[string]int {
	wordCount := make(map[string]int)
	filter := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	text = strings.ToLower(text)
	cleared := filter.ReplaceAllString(text, " ")
	cleared = strings.TrimSpace(cleared) // NOTE: check if i can write like this.
	words := strings.Split(cleared, " ")
	// fmt.Println(words)
	// fmt.Println(len(words))
	wordsCopy := words
	onceTried := make(map[string]int)
	for i := 0; i < len(words); i++ {
		// fmt.Println(i)

		_, exists := onceTried[words[i]]
		if !exists && len(words[i]) > 1 || words[i] == "a" {
			if i < len(words)-1 {
				if words[i+1] == "s" {
					words[i] = strings.Join([]string{words[i], words[i+1]}, "")
					// for j := 0; j < len(words); j++ {
					// 	if words[i] == wordsCopy[j] {
					// 		wordCount[testWord]++
					// 	}
					// }
					//
				}
			}

			for j := 0; j < len(words); j++ {
				if words[i] == wordsCopy[j] {
					wordCount[words[i]]++
				}
			}
		}

		onceTried[words[i]] = 1
	}

	fmt.Println(wordCount)


	return wordCount 
} 
