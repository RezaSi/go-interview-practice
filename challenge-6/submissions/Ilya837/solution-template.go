// Package challenge6 contains the solution for Challenge 6.
package challenge6

import "strings"

func CountWordFrequency(text string) map[string]int {

	del := []string{".", ",", "!", "?", ";", ":", "\t", "\n", "-"}

	for _, d := range del {

		text = strings.ReplaceAll(text, d, " ")
	}

	text = strings.ReplaceAll(text, "'", "")

	split := strings.Split(text, " ")

	res := make(map[string]int)

	for _, word := range split {
		if word != "" {
			res[strings.ToLower(word)]++
		}
	}

	return res
}