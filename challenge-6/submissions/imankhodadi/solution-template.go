package challenge6

import (
	"fmt"
	"strings"
)

func CountWordFrequency(text string) map[string]int {
	words_count := make(map[string]int, 1000) // Pre-allocate for 1000 words
	var builder strings.Builder
	for _, c := range text {
		if c >= 'A' && c <= 'Z' {
			c = c + 32
		}
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			builder.WriteString(string(c))
		} else if c == '\'' {
			continue
		} else {
			if builder.Len() > 0 {
				words_count[string(builder.String())] += 1
				builder.Reset()
			}
		}
	}
	if builder.Len() > 0 {
		words_count[string(builder.String())] += 1
	}
	return words_count
}

func main() {
	fmt.Println(CountWordFrequency("Let's go baby my baby"))
}