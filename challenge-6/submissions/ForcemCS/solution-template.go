// Package challenge6 contains the solution for Challenge 6.
package challenge6


func CountWordFrequency(text string) map[string]int {

	result := make(map[string]int)

	var currentWord []byte

	isAlphanumeric := func(c byte) bool {
		return (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9')
	}

	toLowerChar := func(c byte) byte {
		if c >= 'A' && c <= 'Z' {
			return c + 32
		}

		return c
	}

	for _, char := range text {
	    
	    if char == '\'' {
	        continue
	    }
		if isAlphanumeric(byte(char))  {
			currentWord = append(currentWord, toLowerChar(byte(char)))
		} else if len(currentWord) > 0 {
			word := string(currentWord)
			result[word]++
			currentWord = currentWord[:0]
		}
	}

	if len(currentWord) > 0 {
		word := string(currentWord)
		result[word]++
	}

	return result

}