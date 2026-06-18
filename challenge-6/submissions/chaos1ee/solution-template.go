// Package challenge6 contains the solution for Challenge 6.
package challenge6

import (
	"strings"
	"unicode"
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
	// 将连字符替换为空格，使 "new-lines" 拆分为 "new" 和 "lines"
	text = strings.ReplaceAll(text, "-", " ")
	words := strings.Fields(text)
	freq := make(map[string]int) // map 是关键字，不能作为变量名

	for _, w := range words {
		// 移除所有非字母数字字符（包括中间的标点，如 let's -> lets）
		word := strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsNumber(r) {
				return r
			}
			return -1 // -1 表示从结果中删除该字符
		}, w)
		word = strings.ToLower(word) // 转为小写，实现不区分大小写
		if word == "" {
			continue // 跳过空字符串
		}
		freq[word]++
	}

	return freq
}
