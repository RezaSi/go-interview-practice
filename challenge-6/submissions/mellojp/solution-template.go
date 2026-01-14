package challenge6

import (
	"strings"
	"unicode"
)

func f(r rune) bool {
	return unicode.IsSpace(r) || unicode.Is(unicode.P, r)
}

func CountWordFrequency(text string) map[string]int {
	res := make(map[string]int, 0)
	str := strings.ReplaceAll(text, "'", "")
	array := strings.FieldsFunc(str, f)

	for _, val := range array {
		val = strings.ToLower(val)
		if _, ok := res[val]; !ok {
			res[val] = 1
		} else {
			res[val]++
		}
	}
	return res
}
