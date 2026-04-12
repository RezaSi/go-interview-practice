package regex

import (
	"regexp"
	"unicode"
)

// ExtractEmails extracts all valid email addresses from a text
func ExtractEmails(text string) []string {
	re, err := regexp.Compile(`(?i)[a-z\d._]+(\+[a-z\d]+)?@[a-z\d-]+(\.[a-z]+)+`)
	if err != nil {
		panic("the regexp did not compile")
	}

	result := re.FindAllString(text, -1)

	if result == nil {
		return []string{}
	} else {
		return result
	}
}

// ValidatePhone checks if a string is a valid phone number in format (XXX) XXX-XXXX
func ValidatePhone(phone string) bool {
	re, err := regexp.Compile(`\(\d{3}\) \d{3}-\d{4}`)
	if err != nil {
		panic("the regexp did not compile")
	}

	return len(phone) == 14 && re.MatchString(phone)
}

// MaskCreditCard replaces all but the last 4 digits of a credit card number with "X"
// Example: "1234-5678-9012-3456" -> "XXXX-XXXX-XXXX-3456"
func MaskCreditCard(cardNumber string) string {
	re, err := regexp.Compile(`(\d{4}-?)*\d{4}`)
	if err != nil {
		panic("the regexp did not compile")
	}

	matchIndices := re.FindAllStringIndex(cardNumber, -1)

	if matchIndices == nil {
		return cardNumber
	}

	cardNumRunes := []rune(cardNumber)
	for _, match := range matchIndices {
		for i := match[0]; i < match[1]-4; i++ {
			if unicode.IsDigit(cardNumRunes[i]) {
				cardNumRunes[i] = 'X'
			}
		}
	}

	return string(cardNumRunes)
}

// ParseLogEntry parses a log entry with format:
// "YYYY-MM-DD HH:MM:SS LEVEL Message"
// Returns a map with keys: "date", "time", "level", "message"
func ParseLogEntry(logLine string) map[string]string {
	re, err := regexp.Compile(`(\d{4}-\d{2}-\d{2}) (\d{2}:\d{2}:\d{2}) ([A-Z]+) (.+)`)
	if err != nil {
		panic("the regexp did not compile")
	}

	matchedPieces := re.FindStringSubmatch(logLine)
	if matchedPieces == nil {
		return nil
	}

	return map[string]string{
		"date":    matchedPieces[1],
		"time":    matchedPieces[2],
		"level":   matchedPieces[3],
		"message": matchedPieces[4],
	}
}

// ExtractURLs extracts all valid URLs from a text
func ExtractURLs(text string) []string {
	re, err := regexp.Compile(`https?://(?i)((([a-z\d_]+:.+@)?([a-z\d\-]+)(\.[a-z\d\-]+)+)|localhost)(:\d{4})?(/(?i)[a-z\d.\-_]+)*(\?(?i)[a-z]+=[a-z\d.\-_]+)*(#(?i)[a-z\d]+)?`)
	if err != nil {
		panic("the regexp did not compile")
	}

	matches := re.FindAllString(text, -1)
	if matches == nil {
		return []string{}
	} else {
		return matches
	}
}
